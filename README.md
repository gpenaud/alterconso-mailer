

```

   ______                 __  __                   __  ___      _ __         
  / ____/___ _____ ____  / /_/ /____              /  |/  /___ _(_) /__  _____
 / /   / __ `/ __ `/ _ \/ __/ __/ _ \   ______   / /|_/ / __ `/ / / _ \/ ___/
/ /___/ /_/ / /_/ /  __/ /_/ /_/  __/  /_____/  / /  / / /_/ / / /  __/ /    
\____/\__,_/\__, /\___/\__/\__/\___/           /_/  /_/\__,_/_/_/\___/_/     
           /____/                                                            

```
# Overview

[cagette-mailer 💻](https://github.com/gpenaud/cagette-mailer) is a microservice, developped in python, which allows [cagette-webapp](https://github.com/gpenaud/cagette-webapp) to easily send mail through an arbitrary SMTP server (real or relay). Usage of this microservice is related to [cagette-webapp](https://github.com/gpenaud/cagette-webapp), and should not be used out of this context.

## Requirements

* `curl` (HTTP client for api calls)
* `python` (Python 3.6 up to 3.10 supported)
* `pip` (Python package manager)
* `pipenv` (Python venv manager)
* `Docker`

## Installing

Please clone the repository

```
git clone https://github.com/gpenaud/cagette-mailer.git
```

Go to the repository folder, then execute importer through `pipenv.

**Note**: The api will by default be available on http://127.0.0.1:5000

## Example

Create, Fill then export ENV values from environment.txt by sourcing it:
```
source environment.txt
```

Start cagette-importer in development mode by running:
```
% FLASK_APP=mailer \
  FLASK_ENV=development \
  pipenv run flask run

* Serving Flask app 'mailer' (lazy loading)
* Environment: development
* Debug mode: on
* Running on http://127.0.0.1:5000 (Press CTRL+C to quit)
* Restarting with stat
* Debugger is active!
* Debugger PIN: 137-663-782

```

Then you can query the api by using curl:

```
% curl http://127.0.0.1:5000/healthcheck
  ok
```
