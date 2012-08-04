#!/usr/bin/env rake
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
include ActionView::Helpers::DateHelper

namespace :app do
  desc 'Application status'
  task :status => :environment do
    last = User.last ; puts "-- Users: #{User.count} | #{last.present? ? last.updated_at.to_s : ''}"
    last = Workspace.last ; puts "-- Workspaces: #{Workspace.count} | #{last.present? ? last.updated_at.to_s : ''}"
    last = Grid.last; puts "-- Grids: #{Grid.count} | #{last.present? ? last.updated_at.to_s : ''}"
    last = Row.last; puts "-- Rows: #{Row.count} | #{last.present? ? last.updated_at.to_s : ''}"
    last = Audit.last; puts "-- Audits #{Audit.count} | #{last.present? ? last.updated_at.to_s : ''}"
  end

  desc 'Restarts server'
  task :restart do
    sh 'touch tmp/restart.txt'
  end

  desc 'Grant access for a user to a workspace'
  task :grant_access => :environment do
    puts "Grant access for a user to a workspace"
    if ENV['user'].nil? or ENV['workspace'].nil?
      puts "Usage rake app:grant_access workspace=<workspace uri> user=<user email>"
      return
    else
      user = User.find_by_email ENV['user']
      if user.nil?
        puts "Can't find user"
        return
      end
      workspace = Workspace.find_by_uri ENV['workspace']
      if workspace.nil?
        puts "Can't find workspace"
        return
      end
      puts "Process workspace #{workspace.uuid} for user #{user.uuid}"
    end
  end
  
  desc 'Import data'
  task :import_data => :environment do
    if ENV['file'].nil?
      puts "Usage rake app:import_data file=<file_name>"
      return
    else
      file = ENV['file']
      puts "rake:import_data upload #{file}.xml"
      errors = total_count = total_inserted = total_updated = 0
      Thread.current[:session_locale] = 'en'
      Thread.current[:session_as_of_date] = Date.current
      Thread.current[:session_user_display_name] = User::SYSTEM_ADMINISTRATOR_UUID
      begin
        upload = Upload.create
        upload.upload(file)
        upload.save!
        total_count = total_count + upload.records
        total_inserted = total_inserted + upload.inserted
        total_updated = total_updated + upload.updated
      rescue Exception => invalid
        Entity.log_error "rake:import_data", invalid
        puts "rake:import_data " + invalid.inspect
        errors = errors + 1
      end
      puts "rake:import_data complete " +
                "#{errors} errors, " + 
                "#{total_count} records, " + 
                "#{total_inserted} inserted, " + 
                "#{total_updated} updated."
    end
  end
end