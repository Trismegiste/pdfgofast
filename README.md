# Pdf Go Fast
A ready-copy-paste docker container to generate HTML-to-PDF with a micro-service calling chromium

##What
This is docker config to add to your project (and customize as you wish).

##Where
Just copy-paste the ```docker``` directory into your project root and edit your ```docker-compose.yml``` by appending the content of the ```docker-compose.yml``` from this repo.

##How
This container is using the insanely-optimized ```chromedp/headless-shell``` image that provides a stack for ChromeDriver and Chromium. I've build a http server in Go on top of that.

Calling this microservice is simple as posting a form.
There are two URL provided :
* http://localhost:4444/ : a test form. You provide an html, it returns a PDF. Delete this handler in production. 
* http://localhost:4444/upload : the URL to POST the html file. Only one field in the form named ```file```

See an example in PHP at https://github.com/Trismegiste/eclipse-wiki/blob/master/src/Service/Pdf/ChromiumPdfWriter.php . The method domToPdf posts a form with the Symfony framework to this microservice.

Why posting a form instead of RESTful JSON ? It's lighter, on both side, we're dealing with large binary files.

##Why
I've tried many ways to generate PDF from HTML :
* wkhtml2pdf : fast but old, deprecated and CSS1 only
* calling chromium in a child process with the flag --print-to-pdf : memory intensive, tricky configuration, warning in the log (dbus error, you know what I'm speaking about) and you need to write input before the call.
* API calls to Selenium callling ChromeDriver : nicer, great separation-of-concern than calling chromium but a little bit slower and far more too technological for my needs
* this : K.I.S.S at max level

Of course, this container is the best :) Anyway, it could be better, by using stream, adding https, better log...

##Who
Everyone who need to generate PDF on the backend in a declarative language, therefore HTML/CSS. You can even add Javascript but I don't recommand it, it will be difficult to test and debug. I also recommand you data-url-base64-encode all pictures to prevent any network/DNS problem which can occur depending on your network stack config, and sometimes, Chromium does not wait till all pictures are fully loaded (this could happend when the html page includes many pictures)

