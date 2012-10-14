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
  before_filter :load_credentials, :set_locale, :set_asofdate

  # Defines the application layout.
  layout "application"

  # Includes helpers required for page rendering.
  helper "application", "entity"

  # A feature in Rails that protects against Cross-site Request Forgery (CSRF) attacks.
  # This feature makes all generated forms have a hidden id field.
  # This id field must match the stored id or the form submission is not accepted.
  # This prevents malicious forms on other sites or forms inserted with XSS from submitting to the Rails application.
  protect_from_forgery

protected

  # Loads user credentials.
  def load_credentials
    log_debug "ApplicationController#load_credentials"
    session[:as_of_date] = Date.current if session[:as_of_date].nil?
    Entity.session_as_of_date = session[:as_of_date]
    if user_signed_in?
      Entity.session_user_uuid = current_user.uuid
      Entity.session_user_display_name = current_user
    else
      Entity.session_user_uuid = Entity.session_user_display_name = nil
    end
  end

  # Sets locale based on parameter.
  def set_locale
    log_debug "ApplicationController#set params[:locale]=#{params[:locale]}"
    log_debug "ApplicationController#set I18n.locale=#{I18n.locale}"
    session[:locale] = params[:locale] if params[:locale].present?
    session[:locale] = I18n.default_locale if I18n.locale.nil?
    I18n.locale = session[:locale]
    Entity.session_locale = I18n.locale.to_s
    log_debug "ApplicationController#set Entity.session_locale=#{Entity.session_locale}"
  end
  
  # Sets the as of date based on parameter.
  def set_asofdate
    if params[:as_of_date].present?
      log_debug "ApplicationController#refresh date=#{[:as_of_date]}"
      begin
        requested_date = Date.parse(params[:as_of_date])
      rescue Exception => invalid
        flash[:notice] = t('error.invalid_date', :date => params[:as_of_date])
        redirect_to session[:last_url]
        return
      end
      if requested_date != session[:as_of_date]
        session[:as_of_date] = requested_date
        flash[:notice] = t('general.asofdate', :date => l(session[:as_of_date]))
      end
    end
  end

  # Selects the workspaces available to the connected user.
  def load_workspaces
    log_debug "ApplicationController#load_workspaces"
    @workspaces = Workspace.user_workspaces
  end

  # Keeps track of the current page in the history of navigation.
  def push_history
    session[:history_table] = [] if session[:history_table].nil?
    visited = { :page_title => @page_title, :url => request.url, :when => Time.now }
    session[:prior_page_title] = visited[:page_title]
    session[:prior_url] = visited[:url]
    session[:history_table] << visited if visited[:url] != session[:last_url]
    session[:last_url] = visited[:url]
    found = false
  end

  # Insures the user who signed in is updated with the appropriate user.
  # This method is triggered by Devise after sign in.
  def after_sign_in_path_for(resource_or_scope)
    log_debug "ApplicationController#after_sign_in_path_for(#{resource_or_scope})"
    if user_signed_in?
      current_user.update_user_uuid = current_user.uuid
      current_user.save!
    end
    super
  end

  # Change the as of date session based on the selection of historical data for the given entity.
  def change_as_of_date(entity)
    if entity.begin > session[:as_of_date]
      session[:as_of_date] = entity.begin
      flash[:notice] = t('general.asofdate', :date => l(session[:as_of_date]))
    end
    if entity.end < session[:as_of_date]
      session[:as_of_date] = entity.end
      flash[:notice] = t('general.asofdate', :date => l(session[:as_of_date]))
    end
  end

  # Controller helper used for debug messages.
  def log_debug(message) ; Entity.log_debug(message) ; end

  # Controller helper used for error messages.
  def log_error(message, invalid) ; Entity.log_error(message, invalid) ; end

  # Controller helper used for security messages.
  def log_security_warning(message) ; Entity.log_security_warning(message) ; end
end