from flask import Flask
from flask_restful import Resource, Api, request, reqparse, abort, marshal, fields

from email.message import EmailMessage
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.utils import formatdate

import os
from os import environ

import yaml
import smtplib
import sys

class UnconfiguredEnvironment(Exception):
    """base class for new exception"""
    print("Exception was raised")
    pass

app = Flask(__name__)

@app.route('/healthcheck', methods=['GET'])
def healthcheck():
  return "ok", status.HTTP_200_OK

@app.route('/send', methods=['POST'])
def send_route():
    data = request.get_json(force=True)

    with open(r'config.yaml') as file:
        config = yaml.load(file, Loader=yaml.FullLoader)

    with open(r'secrets.yaml') as file:
        secrets = yaml.load(file, Loader=yaml.FullLoader)

    # environment variables
    smtp_host     = config["smtp_host"]
    smtp_port     = config["smtp_port"]
    smtp_user     = secrets["smtp_user"]
    smtp_password = secrets["smtp_password"]

    if (not smtp_host or not smtp_port or not smtp_user or not smtp_password):
        raise UnconfiguredEnvironment

    server = smtplib.SMTP(smtp_host, smtp_port)
    server.connect(smtp_host, smtp_port)
    server.ehlo()
    server.starttls()
    server.login(smtp_user, smtp_password)

    msg            = MIMEMultipart('alternative')
    msg['Subject'] = data["subject"]
    msg['From']    = data["from_email"]
    recipients     = [ value['email'] for value in data["to"] ]
    msg['To']      = data["from_email"]
    msg['Bcc']     = ', '.join(recipients)
    msg["Date"]    = formatdate(localtime=True)

    part = MIMEText(data["html"], "html")
    for key, header in data["headers"].items():
        part.add_header(key, header)

    msg.attach(part)

    try:
        server.sendmail(msg['From'], recipients, msg.as_string())
        ret = {"status": "sent"}
    except smtplib.SMTPException as e:
        ret = e

    server.quit()
    return ret

app.run(
    host=os.environ.get("MAILER_HOST", "0.0.0.0"),
    port=os.environ.get("MAILER_PORT", "5000"),
)
