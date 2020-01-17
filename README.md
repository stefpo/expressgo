# ExpressGO
An extensible web application server similar to the NodeJS Express server

Similarly to ExpressJS, ExpressGO uses a configurable middleware pipeline, however, due to how GO programs execute compare to NodeJS applications, ExpressGO shows a number of differences compared to its model:

- ExpressGO executes requests concurently. This removes the need to use an asynchronous programming style.
- ExpressGO comes with GoViewEngine as its default view engine. GoViewEngine uses the html/template package. Developers can write their own view engines.
