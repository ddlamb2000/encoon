# encoding: utf-8
#
# Encoon : data structuration, presentation and navigation.
# 
# Copyright (C) 2012 David Lambert
# 
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
# 
# See doc/COPYRIGHT.rdoc for more details.
# Load the rails application
require File.expand_path('../application', __FILE__)

# Initializes the rails application
Encoon::Application.initialize!

# Turns off auto TLS for e-mail.
ActionMailer::Base.smtp_settings[:enable_starttls_auto] = false

APPLICATION_TITLE = "εncooη"

LANGUAGES = {
  "English" => 'en',
  "Français" => 'fr',
  "Español" => 'es'
}