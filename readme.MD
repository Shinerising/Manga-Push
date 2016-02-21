# Manga Push

A PHP program that could automatically push the newest mangas to your kindle devices.

**Attention: The source sites of the mangas are Chinese website so all the mangas are in Chinese. However you can edit the code to make it support more sources but it only pushes Chinese version by now.**

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