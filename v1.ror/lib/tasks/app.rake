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
    Entity.log_debug "Restart server", true
    sh 'touch tmp/restart.txt'
  end

  desc 'Grants access for a user to a workspace'
  task :grant_access => :environment do
    Entity.log_debug "rake:grant_access Grant access for a user to a workspace", true
    if ENV['user'].nil? or ENV['workspace'].nil?
      puts "Usage rake app:grant_access workspace=<workspace uri> user=<user email>"
      return
    else
      user = User.find_by_email ENV['user']
      if user.nil?
        Entity.log_debug "rake:grant_access Can't find user", true
        return
      end
      workspace = Workspace.find_by_uri ENV['workspace']
      if workspace.nil?
        Entity.log_debug "rake:grant_access Can't find workspace", true
        return
      end
      Entity.log_debug "rake:grant_access Process workspace #{workspace.uuid} for user #{user.uuid}", true
      row = WorkspaceSharing.new
      begin
        row.transaction do
          row.create_user_uuid = User::SYSTEM_ADMINISTRATOR_UUID
          row.update_user_uuid = User::SYSTEM_ADMINISTRATOR_UUID
          row.role_uuid = Role::ROLE_TOTAL_CONTROL_UUID
          row.workspace_uuid = workspace.uuid
          row.user_uuid = user.uuid
          row.save!
          saved = true
        end
      rescue ActiveRecord::RecordInvalid => invalid
        Entity.log_debug "rake:grant_access invalid=#{invalid.inspect}", true
        saved = false
      rescue Exception => invalid
        Entity.log_error "rake:grant_access", invalid
        puts "rake:grant_access " + invalid.inspect
        saved = false
      end
    end
    if saved
      Entity.log_debug "Processed workspace #{workspace.uuid} for user #{user.uuid}", true
    end
  end

  desc 'Imports data'
  task :import_data => :environment do
    if ENV['file'].nil?
      puts "Usage rake app:import_data file=<file_name>"
      return
    else
      file = ENV['file']
      Entity.log_debug "rake:import_data upload #{file}.xml", true
      Entity.session_locale = 'en'
      Entity.session_as_of_date = Date.current
      Entity.session_user_uuid = User::SYSTEM_ADMINISTRATOR_UUID
      begin
        upload = Upload.create
        upload.upload(file)
        upload.save!
        Entity.log_debug "rake db:import_data complete " +
                         "#{upload.records} records, " + 
                         "#{upload.inserted} inserted, " + 
                         "#{upload.updated} updated, " +
                         "#{upload.skipped} skipped, " +
                         "#{upload.elapsed} elapsed (ms).", true
      rescue Exception => invalid
        Entity.log_error "rake:import_data", invalid
        puts "rake:import_data " + invalid.inspect
      end
    end
  end
  
  desc 'Updates locales'
  task :locales => :environment do
    Entity.log_debug "rake:locales fetch grids", true
    deleted = 0
    for grid in Grid.all
      grid_uuid = grid.uuid
      grid_name = ""
      for name in grid.grid_locs
        grid_name = name.name if name.locale == 'en'
      end
      has_translation = (grid.has_name or grid.has_description)
      grid_db_table = "rows"
      grid_db_loc_table = has_translation ? "row_locs" : ""
      has_mapping = false
      for mapping in grid.grid_mappings
        grid_db_table = mapping.db_table
        grid_db_loc_table = mapping.db_loc_table if has_translation
        has_mapping = true
      end
      if has_translation and !grid_db_table.blank? and !grid_db_loc_table.blank?
        Entity.log_debug "rake:locales grid #{grid_name} : #{grid_db_table} #{grid_db_loc_table}", true
        collection = Grid.find_by_sql(["SELECT uuid, version " +
                                       "FROM #{grid_db_table} " +
                                       (has_mapping ? "" : "WHERE grid_uuid = :grid_uuid"),
                                       {:grid_uuid => grid_uuid}
                                      ])
        for row in collection
          locales = []
          locs = Grid.find_by_sql(["SELECT id, name, locale " +
                                   "FROM #{grid_db_loc_table} " +
                                   "WHERE uuid = :uuid and version = :version",
                                   {:uuid => row.uuid, :version => row.version}
                                  ])
          for loc in locs
            if (locales.find {|value| loc.locale == value}).present?
              Grid.connection.delete "DELETE FROM #{grid_db_loc_table} " +
                                     "WHERE uuid = '#{row.uuid}' " +
                                     "AND version = #{row.version} " +
                                     "AND locale = '#{loc.locale}' " +
                                     "AND id < #{loc.id}"
              Entity.log_debug "rake:locales delete loc #{loc.locale} #{loc.name}", true
              deleted += 1
            else
              locales << loc.locale
            end
          end
        end
      end
    end
    Entity.log_debug "rake:locales deleted = #{deleted}", true
  end
end