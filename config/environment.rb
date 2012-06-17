# encoding: utf-8
# Load the rails application
require File.expand_path('../application', __FILE__)

# Initializes the rails application
Encoon::Application.initialize!

# Turns off auto TLS for e-mail.
ActionMailer::Base.smtp_settings[:enable_starttls_auto] = false

LANGUAGES = {
  "English" => 'en',
  "Français" => 'fr',
  "Español" => 'es'
}