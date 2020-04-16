

# GoDAM

GoDAM is GoDaddy's proprietary Digital Asset Management tool.

[TOC]



## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and, when it's unusual, how to install them:

1. Windows Server 2016 Standard
   You'll need admin rights

### Set-up Your Dev Environment

#### Install Visual Studio 2017 (Any edition) 

1. Community Edition will work and it's free: 
   https://www.visualstudio.com/thank-you-downloading-visual-studio/?sku=Community&rel=15
   1. During installation, under **Workloads**, select 
      1. **ASP.NET and Web Development**
      2. **Data Storage and Processing** 
      3. **Node.js Development** 
      4. **.NET Core  cross-platform development**

#### Get Git

Most machines already have Git installed. 

1. Check if you have Git.  Open a command prompt and simply type `git`.  If it's not installed, you can download the client here: https://git-scm.com/downloads

#### Clone the repository

Cloning the repository is part of our GitFlow branching model and the instructions can be found under the [Get Started With the Repo](#Get-Started-With-the-Repo) section below.  *But come back here after you've done that!*

[This section describes how our team collaborates with GitHub](#This-is-How-We-Do-It).  ***Read it. Bookmark It. Live It.***

#### Enable Docker support

1. Open an elevated PowerShell session and run the following commands.

   First, install the Docker-Microsoft PackageManagement Provider from the PowerShell Gallery.

   ```
   Install-Module -Name DockerMsftProvider -Repository PSGallery -Force
   ```

   Next, you use the PackageManagement PowerShell module to install the latest version of Docker.

   ```
   Install-Package -Name docker -ProviderName DockerMsftProvider
   ```

1. Reboot

   ```
   restart-computer -Confirm
   ```



#### Install and Configure MySQL Server
1. Install https://dev.mysql.com/downloads/installer/

2. Create a new user godam with password 12qw!@QW

3. Assign the user the DBA priviledge.

4. Open MySQL WorkBench to connect to your local instance, and run the script  \dam\Source\DB\GoDam\MySQLCreateSchema.mysql to create the database


#### Configure the Visual Studio project

You should have your own appsettings.Development.json file which specify your environment.  You can use following as the template for your appsettings.Development.json, and put it in Source\WebApplications\GoDAMWeb. If you do not understand a setting’s meaning, do not change that.

Use the json file attached in email
​

## Running GoDAMWeb Locally

### Using IIS Express

You should set **GoDAMWEB** as the start project, and run against your local IIS Express server.  If you do that, it should work for OKTA/GoDaddy login automatically.

### Using IIS and SSO

1. Make sure your application's Azure AD settings are correct.  

   ```json
    "AzureAd": {
      "Instance": "https://login.microsoftonline.com/",
      "Domain": "secureservernet.onmicrosoft.com",
      "TenantId": "d5f1622b-14a3-45a6-b069-003f8dc4851f",
      "ClientId": "fb952138-bd16-4e52-bace-f1af7b770e5a", //this is client ID for local develoment env http://localhost:44352, please config your app to run on that port on local IIS server
      "CallbackPath": "/Account/LoginDone",
      "CallbackHost": "http://localhost:44352"
    },
    ```
2. Turn off JavaScript debugging on Chrome (a new feature introduced in VS 2017 that throws an error every time you run a debugging session). Go to **Tools> Options > Debugging > General** and turn off the setting for **Enable JavaScript Debugging for ASP****.NET(Chrome and IE)**.

3. Enable IIS

   1. Open **Server Manager**.
   2. Under **Manage** menu, select **Add Roles and Features**
   3. Select **Role-based or Feature-based Installation**
   4. Select the appropriate server (local is selected by default)
   5. Select **Web Server (IIS)**
   6. No additional features are needed for IIS, so click **Next**
   7. Click **Next**
   8. Customize your installation of IIS, or accept the default settings that have already been selected for you, and then click **Next**
   9. Click **Install**
   10. When the IIS installation completes, the wizard reflects the installation status
   11. Click Close to exit the wizard.

4. Add Development time IIS support to Visual Studio

   Launch the Visual Studio installer to modify your existing Visual Studio installation. In the installer select the **Development time IIS support** component which is listed as optional component under the **ASP****.NET and web development** workload. This will install the ASP.NET Core Module which is a native IIS module required to run ASP.NET Core applications on IIS

5. ASP.NET Core Module (ANCM) lets you run ASP.NET Core applications behind IIS, using IIS for what it's good at (security, manageability, and lots more) and using [Kestrel](https://docs.microsoft.com/en-us/aspnet/core/fundamentals/servers/kestrel) for what it's good at (being really fast), and getting the benefits from both technologies at once.

   Download the Windows (Server Hosting) installer and run the exe from an Administrator command prompt:

   https://aka.ms/dotnetcore-2-windowshosting

 (**Ignore step 6 to step 8  at this moment, you do not need it yet**)
 
6. Install RabbitMQ 
   1. Install Erland from http://www.erlang.org/downloads latest windows version
   2. Install Rabbit MQ from http://www.rabbitmq.com/install-windows.html
   3. Insatll Rabbit MQ Management Plugin https://www.rabbitmq.com/management.html 
   Run in the the command line (Run as administrator) (cd C:\Program Files\RabbitMQ Server\rabbitmq_server-3.7.14\sbin)
   ```
   rabbitmq-plugins enable rabbitmq_management
   ```
   4. Add Environment Variables ERLANG_HOME (C:\Program Files\erl10.3) if not existing. 
   5. Add one more line (%ERLANG_HOME%\bin) in PATH Environment Variables 
   6. Restart the RabbitMQ service
   7. Web UI is at http://localhost:15672/
      


8. Check your work

   1. Surf to https://gddam.dev-gdcorp.tools/  The browser might complain "Your connection is not private".  Just move forward by clicking **Proceed to godam-dev.godaddy.com (unsafe)**
   2. You should see the project's home page
   3. If that works, you should be able to sign in using GoDaddy's SSO (Okta) by clicking **Sign In** or any of the pages that require authorization.

9. You're Done setting up your dev machine!  

## Debugging GoDamWeb

The Visual Studio Tools for Docker provides a consistent way to develop in and validate your application locally in a Linux Docker container. You don't have to restart the container each time you make a code change.  This section illustrates how to use the "Edit and Refresh" feature to start GoDAMWeb in a local Docker container, make any necessary changes, and then refresh the browser to see those changes. This section also shows you how to set breakpoints for debugging.






On GitHub (server side), we have two main branches:

**Master branch** is only used for product release. Anything checked-in to that branch will be PROD code.  Code will need go through a PULL request to get into the Master branch; NO direct check-ins to Master branch. The build from master branch can be deployed to PROD at any time.  This ensures auto-rollback. All Master check-ins should have a TAG associated which indicates the build number. 

**Develop branch** is the SHARED and default development branch. All tested code can be merged into development branch and should go through PULL request instead of direct push.

Other branches are the individual's feature/development branches; each contributor manages independently.  The best practice is to create a branch with your own alias or the feature that you are working on(maybe a JIRA ticket number).

### Get Started With the Repo

1. At the command prompt, create a folder where you want to store your repository and `cd` into it.

2. Clone the repository 
   `git clone https://github.com/GD-China-DAM/dam.git`

3. Change directory to the dam repo 
   `cd dam`

4. Switch the branch on your local machine
   `git checkout development` 

5. `git branch`  
   You should see your develop branch has a star next to it, which means you are in that branch.

6. Ensure that you have the latest code from develop branch
   `git pull`  

7. Create your own working branch
   `git branch -c "myfeature"` 

8. Switch you to your own branch, (myfeature branch), on your local machine

   `git checkout myfeature` 

9. Now you can work on your own branch freely. Use `git commit –m "My brief description of this commit"` such that you can commit/rewind your changes in your local repository.

### This is How We Do It

*Please READ below section carefully.*

So you've worked hard <u>in your branch</u>, your feature is good to go and you want to get your code into the development branch.  Great!  Here's what you do:

#### 1. Code Review 

1. `git branch` 
   You should see your your own branch has a star next to it.  That means your branch is checked out and your changes have been made in that branch.
2. Ensure you have checked in everything you want to check in 
   `git commit`
3. `git push` so your remote feature branch is updated with your change.
4. Create a pull request so you can get feedback from others (note: here the pull request is only for code review, it’s not for merge)
   1. Surf to https://github.secureserver.net/gdi/dam/ 
   2. Click **New Pull request**
   3. Select **development** (left side), and your feature branch (right side)
   4. Click **Create pull request**
      ![](README.md.images/NewPullRequest.png)
5. Fix the issues that others raised in the pull request, and go back to step 1 to start another pull request, until everything is fixed.  
6. Now you're ready to merge

#### 2. Merge

Did you complete your code review?  No?  Do it and come back when its done.

7. Check out the Develop branch
   `git checkout develop`

8. Get the latest code from remote development branch
   `git pull`

9. Merge your code into the development branch on your local machine
   `git merge --squash myfeature` 

10. Resolve any conflicts with visual studio's Editor,

11. Since resolving conflict may take a while, and at this moment, the code in remote repo has changed , you will need to pull code again. Here is the way to do it Run 
     `git add .` to add all your changed files again
     `git stash` to put your work into a temperary workspace to allow git pull
     'git pull'  this will garentee success, since your work has been put into a seperate workspace, so no conflicts at all
     'git stash pop' to reapply your work into local develop branch, this may cause conflict again, and you will need to resolve it
 
12.  Continue step11, until no conflict after your "git stash pop'    
   
    until all conflicts are resolved. Now, with your work merged into your updated local development branch and it still works, you're ready to push your local development branch.    
     
13. `git commit -m "ticketid" -a`  (for example `git commit -m "DAM-1234" -a` )

this will make ONE commit with name "ticketid" into your local development branch

14. In principle, each commit into development branch should come with a ticket in the commit message. The only exception is that the fix is really minor, and you can describle your fix in couple of words such as "typo fix for word ordr" 

#### 3. Push Develop Branch

1. Make sure you're working with your local development branch
   `git status' to ensure you are in development branch
2. push your local development branch to remote repository
   `git push`
3. Delete your feature branch if it's no longer needed at https://github.com/GD-China-DAM/dam

### Git Command Quick Reference

#### Create a new branch

```
git -b mybranchname
```

#### Switch between branches

```
git checkout mybranchname
```

#### Commit change to local repository  

```
git commit –m "message"
```

#### Check which files you want to commit:  

```
git status 
```

It's a good idea to do this before every commit.

#### Stage all changed files to your staging area “

```
git stage *
```

Use this to see what files you are tracking at this moment.

#### Download the latest stuff from remote repository 

```
git pull
```

#### Push files to remote repository 

```
git push
```

#### Undo one file in your local change and in your staging area 

```
git checkout mybranchname --myfilename
```

That undoes changes in that file that are not committed.

#### Undo one file in your local change 

```
git checkout --filename
```

That undoes changes in that file that are not staged = undo any thing changed after latest `git stage`. 

#### Undo all unstaged files

```
git checkout .
```

This replaces all local files with the last snapshot you created through git stage.  This is a risky command.

#### Revert all your local changes

```
git reset --hard  
```

This discards anything you have changed on your local machine, and revert them (both in the staging area and local) back to the latest committed version in the branch.  Do I really need to tell you this is a risky command?

## MySQL Database Update and Model Rebuilt
```
Scaffold-DbContext "server=localhost;port=3306;user=godam;password=12qw!@QW;database=GoDAM;" "Pomelo.EntityFrameworkCore.MySql" -OutputDir Models -force -verbose
```

## MySQL Database Backup and Restore

### Script to backup the database: (This will great a dump file such as godamdump2018_08_01.sql.gz)
```
 /usr/local/mysql/bin/mysqldump -f -Q -R --opt --skip-lock-tables -u godam -p.... --databases godam | gzip > /home/szhang/mysqlbackup/godamdump$(date+"%Y_%m_%d").sql.gz
```
 
### Script to unzip the backup:
```
gzip -d godamdump2018_08_16.sql.gz
```

### Restore the database:
```
/usr/local/mysql/bin/mysql -u godam -pPassword godam < file.sql
``` 


# Note
1) CDN to prevserver: If we want to clean up ceph server, we will need to republish this URL http://img2.wsimg.com/cdn/Email/v1/All/17f947a2-9066-45a7-b3fd-e94b524ce966/GD_WORDMARK_RGB_BLACK.png

# Setup Local Debugging for SSO

## 1. Disable HTTPS2 on your dev machine
1) Start the Windows Registry Editor.
2) Navigate to the registry key HKEY_LOCAL_MACHINE\System\CurrentControlSet\Services\HTTP\Parameters.
3) Add 2 new REG_DWORD values, EnableHttp2Tls and EnableHttp2Cleartext, to this registry key.
4) Set both values to 0.
5) Reboot machine

## 2. Enable Chrome to trust local certifcate
1) Launch Chrome
2) Copy/Paste the following string to Address bar chrome://flags/#allow-insecure-localhost
3) enable "Allow invalid certificates for resources loaded from localhost."
4) Restart Chrome

## 3 Optional- re-arrange the Crypto sequence on your machine
1) If you see have issue to connect through Chrome
2) Visit https://www.nartac.com/Products/IISCrypto  , download the app, and run it
3) Click "Best Practice" in the left bar, that will disable couple crypto kit on your server
4) Reboot your machine

# Setup you local dev env for SSO
1) Get the latest appsettings.Development.json from the Source\Configuration folder
2) It contains following entry

 "GoDaddySSO": {
    "Name": "GoDaddySSO",
    "SigninUrl": "https://sso.dev-gdcorp.tools?realm=jomax&app=gddam",
    "SignoutUrl": "https://sso.dev-gdcorp.tools/logout?realm=jomax&app=gddam",
    "SSOUrl": "https://sso.dev-gdcorp.tools",
    "Enabled": true
  },

3) Enabled value allows you to switch between SSO and OKTA experience, once we fully tested SSO, we will disable OKTA in the future, for now you can use Enabled : false to continue delveop your feature if your dev work is blocked somehow by SSO

4) You will need to add following value into your host file

127.0.0.1 gddam.dev-gdcorp.tools

5) When you debug your application, select IIS Express SSO configure in the Visual Studio, that will launch the HTTPS version with SSO on dev-gdcorp.tools domain instead of godaddy.com domain


