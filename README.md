# Privates Lauftagebuch – Private Running Journal

The private running journal is a modern desktop app used to track one's running training.

The project is currently in a beta state.

## Features

* __Design Philosophy__: All data is only stored on the desktop itself. This is what sets this running journal apart from all the other similar apps in the web.
* All data is stored as readable files in common formats. Thus, even without the app, the running journal's data is readable.
* Creating and updating tracks
* Creating journal entries after the training
* Customizable OSM tile server

## Planned Features

* Analytics of trainings
* a lot of UI improvements
* Windows and cross-platform copying of journal data (some "/" might be hardcoded at the moment ;)
* Import and export of tracks via the UI (due to the open design principle, it's already possible to simply access the file system.)
* Improved error handling (many errors are ignored in the current beta version)
* many more …

## Contributing

This project is built using Wails and VueJS. In order to contribute code to this project, you need to setup Wails and pnpm on your machine.

If Wails and pnpm are installed, the app can be started in development mode using `wails dev -appargs "appdata"`.

The `appargs` argument specifies the working directory of the application. If ommited, it will use use `~/.private-running-journal`.

If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115.

To build the app, use `wails build`.

This is the official Wails Vue-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

Of course, you can also contribute by giving your ideas and feedback to the project, preferably via the Github page.
