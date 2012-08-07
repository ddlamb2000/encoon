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
errors = total_count = total_inserted = total_updated = total_skipped = total_elapsed = 0
for file in ["administrator", "credentials", "system", "kinds", "roles", "countries"]
  Entity.log_debug "rake db:seed upload #{file}.xml", true
  Thread.current[:session_locale] = 'en'
  Thread.current[:session_as_of_date] = Date.current
  Thread.current[:session_user_display_name] = User::SYSTEM_ADMINISTRATOR_UUID
  begin
    upload = Upload.create
    upload.upload("db/#{file}.xml")
    upload.save!
    Entity.log_debug "rake db:seed complete " +
                     "#{upload.records} records, " + 
                     "#{upload.inserted} inserted, " + 
                     "#{upload.updated} updated, " +
                     "#{upload.skipped} skipped, " +
                     "#{upload.elapsed} elapsed (ms).", true
    total_count = total_count + upload.records
    total_inserted = total_inserted + upload.inserted
    total_updated = total_updated + upload.updated
    total_skipped = total_skipped + upload.skipped
    total_elapsed = total_elapsed + upload.elapsed
  rescue Exception => invalid
    Entity.log_error "rake db:seed", invalid
    puts "rake db:seed " + invalid.inspect
    errors = errors + 1
  end
end
Entity.log_debug "rake db:seed total " +
                 "#{errors} errors, " + 
                 "#{total_count} records, " + 
                 "#{total_inserted} inserted, " + 
                 "#{total_updated} updated, " +
                 "#{total_skipped} skipped, " +
                 "#{total_elapsed} elapsed (ms).", true
                 