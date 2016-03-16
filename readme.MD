# Manga Push

A PHP program that could automatically push the newest mangas to your kindle devices.

**The resource manga websites we used have disabled the images copy function, so that this project have to close:( I will keep seeking new resources and thank you for your attention.**

## Features

* Automatically trace new mangas by RSS feed 
* Intelligently fetch & compress images and create PDF
* Provide RSS feed for further requirement
* You can check the DEMO here: [Manga Push](http://pku-demo-ebooks.lionfree.net/manga/)

## How to use

1. Edit the config for your personal usage
2. Add user and manga list in settings.xml
3. Upload all the files to your server
4. Add a feed service (like IFTTT) to regularly trigger the program, the feed url is ` [YOUR_SERVER_HOST]/feed.php?user=[USER_NAME] `

## Third Party Libraries

* Fpdf www.fpdf.org
* PHPMailer https://github.com/PHPMailer/PHPMailer