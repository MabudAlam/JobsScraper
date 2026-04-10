#!/usr/bin/env python3
"""
Validate company slugs against Ashby API.
Removes invalid slugs from companies.json.

Usage:
    python scripts/validate_companies.py

The script will:
1. Load companies from companies.json
2. Validate each slug via Ashby GraphQL API
3. Remove invalid companies
4. Save the updated companies.json
"""

import json
import urllib.request
import urllib.error
import time
import sys
import os


def validate_slug(slug):
    """Call Ashby API to validate if slug exists"""
    query = """query ApiOrganizationFromHostedJobsPageName($organizationHostedJobsPageName: String!, $searchContext: OrganizationSearchContext) {
  organization: organizationFromHostedJobsPageName(
    organizationHostedJobsPageName: $organizationHostedJobsPageName
    searchContext: $searchContext
  ) {
    ...OrganizationParts
    __typename
  }
}

fragment OrganizationParts on Organization {
  name
  publicWebsite
  hostedJobsPageSlug
  timezone
  __typename
}"""

    payload = json.dumps(
        {
            "operationName": "ApiOrganizationFromHostedJobsPageName",
            "variables": {
                "organizationHostedJobsPageName": slug,
                "searchContext": "JobBoard",
            },
            "query": query,
        }
    ).encode("utf-8")

    req = urllib.request.Request(
        "https://jobs.ashbyhq.com/api/non-user-graphql?op=ApiOrganizationFromHostedJobsPageName",
        data=payload,
        headers={
            "Content-Type": "application/json",
            "Accept": "*/*",
            "Accept-Language": "en-US,en;q=0.9",
            "apollographql-client-name": "frontend_non_user",
            "apollographql-client-version": "0.1.0",
            "origin": "https://jobs.ashbyhq.com",
            "user-agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Mobile Safari/537.36",
        },
        method="POST",
    )

    try:
        with urllib.request.urlopen(req, timeout=10) as response:
            data = json.loads(response.read().decode("utf-8"))
            if (
                "data" in data
                and "organization" in data["data"]
                and data["data"]["organization"]
            ):
                return True, data["data"]["organization"].get("name", "")
            return False, ""
    except urllib.error.HTTPError as e:
        return False, f"HTTP {e.code}"
    except Exception as e:
        return False, str(e)


def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(script_dir)
    companies_file = os.path.join(project_root, "companies.json")

    with open(companies_file, "r") as f:
        companies = json.load(f)

    print(f"Validating {len(companies)} companies...")
    print("-" * 50)

    valid = []
    invalid = []
    errors = []

    for i, company in enumerate(companies):
        slug = company["AshbySlug"]
        name = company["Company"]

        valid_result, org_name = validate_slug(slug)

        if valid_result:
            valid.append(company)
            status = f"✓ VALID ({org_name})"
        else:
            invalid.append(company)
            status = f"✗ INVALID ({org_name})"
            errors.append((name, slug, org_name))

        print(f"[{i + 1}/{len(companies)}] {name} ({slug}): {status}")

        time.sleep(0.1)

    print("-" * 50)
    print(f"\nResults:")
    print(f"  Valid:   {len(valid)}")
    print(f"  Invalid: {len(invalid)}")

    if invalid:
        print(f"\nInvalid companies removed:")
        for name, slug, err in errors:
            print(f"  - {name} ({slug})")

        with open(companies_file, "w") as f:
            json.dump(valid, f, indent=2)

        print(f"\nUpdated companies.json with {len(valid)} valid companies")
    else:
        print("\n✓ All companies are valid!")


if __name__ == "__main__":
    main()
