# cagette-mailer
A tiny micro-service  mailer for cagette project

## execute mailer for development environment
  SMTP_HOST= \
  SMTP_PORT= \
  SMTP_USER= \
  SMTP_PASSWORD= \
  FLASK_APP=mailer.py \
  FLASK_ENV=development \
flask run -h 0.0.0.0 -p 5000
