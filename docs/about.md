# About the project

## File structure

```bash
.
├── app # api section
│   ├── app
│   ├── app.go # main()
│   ├── mysqldb.go # mysql database utility
│   ├── res # media resources used for streaming
│   │   └── image
│   │       └── A.jpg
│   ├── routes.json
│   ├── sample.json
│   ├── static
│   │   └── index.html
│   ├── status.json
│   └── util.go # library for utlity
├── docs
│   └── about.md
├── launch # backend launcher
│   ├── build.sh # build the APP binary out of source code
│   ├── start.sh # run the APP binary out of source code
│   └── test.sh # merge build and start
├── LICENSE
└── README.md
```

## Some settings

### To remove colory logger

I've just implemented the logger of my choice (I like designs) which you can change by just below steps.

1. `remove` or `comment` switch>>case for choosing logger format. (line 113 to 180)
2. `remove` or `comment` color variables declaration.

You can keep the normal logger by extracting out of `default:` block.

### To disable documenting routes in routes.json

`remove` or `comment` route recorder block; line 192 to 198.

## To use this boilerplate on your own backend system

1. Clone the repository. It is recommended to clone out out existing source code as you might have to face version control issues.

2. From `app.go`, get the `var()` and `main()` function loaded. `import()` will get loaded on its own and if not, you have to do it manually depending upon the IDE you are using. (Or copy entire `app.go` into your project), but you must copy `status.json` into your source directory or build up its replicate.

3. Now, set the logger and routes at your preferences. For static server's need, set the correct route and path at `e.Static()` section [line 203]

4. Now, for REST handler to send JSON response check block [line 190 - 193]. For connecting to MySQL database, copy `mysql.go` file into your project. Study the pipeline for method: `sendSampleDBData()` and for local JSON response, copy `util.go` into your project and study the pipeline for `sendSampleLocalData()` method.

5. For streaming file and other resources, study `sendImage()` method, which sends down the file as per the request is validation (you can remove the validation as it is just there for demonstration purposes).

## About using Docker

You are free to use containerization in your own project and the codebase of this broilerplate doesnot obstruct you from doing so.

Just have your APP ENTRYPOINT well mentioned and you're good to setup docker image from your Dockerfile.

Refer the following blog for simple Docker setup: [Using Docker](https://medium.com/@yuiltripathee/using-docker-for-the-first-time-fd804ae6c209)
