import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.utils import formatdate
import os
from os import environ

from flask import Flask
from flask_restful import Resource, Api, request, reqparse, abort, marshal, fields

class UnconfiguredEnvironment(Exception):
    """base class for new exception"""
    pass

app = Flask(__name__)

@app.route('/send', methods=['POST'])
def send_route():
    data = request.get_json(force=True)

    # environment variables
    smtp_host     = os.environ.get("SMTP_HOST", None)
    smtp_port     = os.environ.get("SMTP_PORT", None)
    smtp_user     = os.environ.get("SMTP_USER", None)
    smtp_password = os.environ.get("SMTP_PASSWORD", None)

    if (not smtp_host or not smtp_port or not smtp_user or not smtp_password):
        raise UnconfiguredEnvironment

    server = smtplib.SMTP(smtp_host, smtp_port)
    server.set_debuglevel(1)

    print(smtp_host)
    print(smtp_port)
    print(smtp_user)
    print(smtp_password)

    server.connect(smtp_host, smtp_port)
    server.ehlo()
    server.starttls()
    server.login(smtp_user, smtp_password)

    msg            = MIMEMultipart('alternative')
    msg['Subject'] = data["subject"]
    msg['From']    = data["from_email"]
    recipients     = [ value['email'] for value in data["to"] ]
    msg['To']      = ', '.join(recipients)
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
    host=os.environ.get("MAILER_HOST", "127.0.0.1"),
    port=os.environ.get("MAILER_PORT", "5000"),
)
