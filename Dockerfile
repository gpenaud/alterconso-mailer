FROM python:3.8-slim-buster

RUN pip3 install pipenv

ENV PROJECT_DIR /usr/src/flaskbookapi

WORKDIR ${PROJECT_DIR}

COPY Pipfile .
RUN pipenv install --deploy --ignore-pipfile

COPY . .

EXPOSE 5000
CMD ["pipenv", "run", "python", "mailer.py"]
