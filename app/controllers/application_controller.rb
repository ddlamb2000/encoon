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
class ApplicationController < ActionController::Base
  before_filter :load_credentials

  # Application layout
  layout "application"

  # Include all helpers, all the time
  helper :all 

  # Security
  protect_from_forgery

protected

  def load_credentials
    session[:as_of_date] = Date.current if session[:as_of_date].nil?
    Entity.session_as_of_date=session[:as_of_date]
    if user_signed_in?
      Entity.session_user_uuid = current_user.uuid
      Entity.session_user_display_name = current_user
    else
      Entity.session_user_uuid = nil
      Entity.session_user_display_name = nil
    end
    session[:locale] = params[:locale] if params[:locale]
    I18n.locale = session[:locale] || I18n.default_locale
    Entity.session_locale = I18n.locale.to_s
  end

  def load_workspaces
    log_debug "ApplicationController#load_workspaces: " + 
                "user_uuid=#{Entity.session_user_uuid}?"
    @workspaces = Grid.select_entity_by_uuid(Grid, Workspace::ROOT_UUID)
    @workspaces.load_cached_grid_structure if @workspaces.present?
  end

  def push_history
    session[:breadcrumb_table] = [] if session[:breadcrumb_table].nil?
    session[:history_table] = [] if session[:history_table].nil?

    visited = {
                :page_title => @page_title, 
                :url => request.url, 
                :page_icon => @page_icon, 
                :when => Time.now
              }

    session[:prior_page_title] = visited[:page_title]
    session[:prior_url] = visited[:url]
    session[:prior_page_icon] = visited[:page_icon]
    
    if visited[:url] != session[:last_url] 
      session[:history_table] << visited
    end

    session[:last_url] = visited[:url]

    found = false
    breadcrumbs = []
    session[:breadcrumb_table].collect do |item|
     if request.url == item[:url]
        found = true
     end
     if !found
        breadcrumbs << item
        session[:prior_page_title] = item[:page_title]
        session[:prior_url] = item[:url]
        session[:prior_page_icon] = item[:page_icon]
     end
    end
    breadcrumbs << visited
    session[:breadcrumb_table] = breadcrumbs
  end
  
  def change_as_of_date(entity)
    if entity.begin > session[:as_of_date]
      session[:as_of_date] = entity.begin 
      flash[:notice] = 
        "As of date changed to #{session[:as_of_date].to_s}" 
    end
    if entity.end < session[:as_of_date]
      session[:as_of_date] = entity.end 
      flash[:notice] = 
        "As of date changed to #{session[:as_of_date].to_s}" 
    end
  end
  
  def lock_as_of_date
    session[:lock_as_of_date] = true 
  end
  
  def unlock_as_of_date
    session[:lock_as_of_date] = false 
  end
  
  def param_begin_date
    if params[:begin_date].present?
      Date.parse params[:begin_date]
    else
      Entity.begin_of_time
    end
  end

  def log_debug(message)
    Entity.log_debug(message)
  end

  def log_error(message, invalid)
    Entity.log_error(message, invalid)
  end

  def as_of_date_clause(synonym)
    Entity.as_of_date_clause(synonym) 
  end
end