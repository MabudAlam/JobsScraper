# Jobs Scraper
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

![Screenshot](https://github.com/Pavel401/Jobs-Scraper/assets/47685150/bec93412-7c81-443b-990a-9fc0d92bfadc)

## Table of Contents

- [Introduction](#introduction)
- [Mobile App](#mobile-app)  
- [Prerequisites](#prerequisites)
- [Endpoints](#endpoints)
- [Contributing](#contributing)
- [License](#license)
- [Contributors](#contributors)
## Introduction

The Jobs Scraper is a powerful Go application designed for scraping job postings from a variety of websites and storing them in an SQLite database. Additionally, it offers endpoints for retrieving and synchronizing job postings, making it a valuable tool for job seekers and developers alike.

## Mobile App

Repo: https://github.com/Pavel401/JobScraper-Mobile

## Prerequisites

Before getting started with the Jobs Scraper, please ensure you have the following prerequisites installed:

- [Go Programming Language](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)



## How to setup

To run the application, you have to create a `.env` file. To configure you can set environment variables in the file. For example:

```env
SYNC_WITH_SQL_PASSWORD=password
```


```bash
go run main.go
```

## Endpoints

The Jobs Scraper package provides a set of powerful endpoints for scraping, managing, and synchronizing job postings. Here's a brief overview of the available endpoints:

- `/atlassian`, `/amazon`, `/coursera`, `/freshworks`, `/gojek`, `/mpl`: Scrapes job postings from specific websites.
- `/syncwithSql`: Synchronizes job postings with a SQL database.
- `/getallJobsFromSQL`: Retrieves job postings from a SQL database.
- `/`: Serves an HTML file (static/base.html).


## Contributing

Contributions to this project are highly encouraged! If you have ideas for improvements or come across any issues, please don't hesitate to open an issue or submit a pull request on the [GitHub repository](https://github.com/Pavel401/Jobs-Scraper).


This documentation aims to simplify and guide the way beginners make their first contribution. If you are looking to make your first contribution, follow the steps below.

#### If you don't have git on your machine install it.

## Fork this repository

 Fork this repository by clicking on the fork button on the top of this page.
 This will create a copy of this repository in your account.


 <img src="assests\readme_imgs\FORK.png"></img>


## Open in Github Desktop

Click on the green icon indicating "code" and open your Github desktop.

<img src="assests\readme_imgs\open in Git Desktop.png"></img>


## Clone the repository

 <img src="assests\readme_imgs\clone.png"></img>

 Now clone the forked repository to your machine. Go to your GitHub account, open the forked repository, click on the code button and then click on open with Github Desktop ot you can click on  _copy to clipboard_ icon if you want to use git bash.

**Note: Further command are for Git bash users not for the Git Desktop Users**

> For Git Desktop users ---> Now just click on Open with VS code and start your contribution.

 Example:
 <img src="assests\readme_imgs\Vscode img.png"></img>


## Commit Changes

 After you have updated the files,click on 'Commit to main' and then click on push origin.

 Example :
 <img src="assests\readme_imgs\commit.png"></img>


 Now Come back to Github web and click on contribute to submit your changes for review.


***
## Steps for Git Bash Users 


Open a terminal and run the following git command:

```
git clone "url you just copied"
```

where "url you just copied" (without the quotation marks) is the url to this repository (your fork of this project). See the previous steps to obtain the url.

For example:


```
git clone https://github.com/this-is-you/fossc.git
```


 Here you're copying the contents of the first-contributions repository on GitHub to your computer.

## Create a branch

Change to the repository directory on your computer (if you are not already there):

```
cd desktop
```

Now create a branch using the `git checkout` command:

```
git checkout -b <your-new-branch-name>
```
For example:

```
git checkout -b add-new-file
```

(The name of the branch does not need to have the word _add_ in it, but it's a reasonable thing to include because the purpose of this branch is to add your name to a list.)

## Make necessary changes and commit those changes

Now open add or edit file in a text editor. Add code for any existing algorithm in other language or add some new algorithms. Make sure to update correspond README.md file if needed. Now, save the file.

If you go to the project directory and execute the command `git status`, you'll see there are changes.

Add those changes to the branch you just created using the `git add` command:

```
git add "name of the file you add or edit"
```

Now commit those changes using the `git commit` command:

```
git commit -m "Add message for the change"
```

## Push changes to GitHub

Push your changes using the command `git push`:

```
git push origin <add-your-branch-name>
```

replacing `<add-your-branch-name>` with the name of the branch you created earlier.

## Submit your changes for review

If you go to your repository on GitHub, you'll see a `Contribute` button. Click on that button.

<img src="assests\readme_imgs\Contribute.jpeg"></img>

click on `Open pull request`.

<img src="assests\readme_imgs\open pull request.png"></img>

click on `Create pull request`.

<img src="assests\readme_imgs\create request.png"></img>


click on `Create pull request`.

## Where to go from here?

Congrats! You just completed the standard _fork -> clone -> edit -> pull request_ workflow that you'll encounter often as a contributor!


## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/debayangg"><img src="https://avatars.githubusercontent.com/u/66942246?v=4?s=100" width="100px;" alt="Debayan Ghosh"/><br /><sub><b>Debayan Ghosh</b></sub></a><br /><a href="https://github.com/Pavel401/Jobs-Scraper/commits?author=debayangg" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->


## License

This project is licensed under the [GPL-3.0 license] License - see the [LICENSE](LICENSE) file for details.

