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

1. `remove` or `comment` switch>>case for choosing logger format.
2. `remove` or `comment` color variables declaration.

You can keep the normal logger by extracting out of `default:` block.
