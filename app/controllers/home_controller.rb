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
class HomeController < ApplicationController

  def index
    @page_title = ""
    @page_icon = "home"
    unlock_as_of_date
    push_history
    respond_to do |format|
      format.html
    end
  end
  
  def refresh
    log_debug "HomeController#refresh date=#{params[:home][:session_date]}"
    session[:as_of_date] = Date.strptime(params[:home][:session_date], I18n.t('datepicker.decode'))
    redirect_to session[:last_url]
  end

  def history
    respond_to do |format|
      format.js
    end
  end
  
  def hide_history
    respond_to do |format|
      format.js
    end
  end

  def import
    log_debug "HomeController#import"
    @page_title = "Import Data"
    @page_icon = "import"
    @upload = Upload.new
  end
  
  def upload
    log_debug "HomeController#upload"
    @upload = Upload.new
    @upload.create_user_uuid = session[:user_uuid]
    @upload.update_user_uuid = session[:user_uuid]
    respond_to do |format|
      if @upload.update_attributes(params[:upload])
        flash[:notice] = "File uploaded."
        format.html { redirect_to :action => "index" }
      else
        format.html { render :action => "import_data" }
      end
    end
  end

  def export_system
    log_debug "HomeController#export_system"
    @workspace = Workspace.select_entity_by_uuid(Workspace, Grid::SYSTEM_WORKSPACE_UUID)
    if @workspace.present?
      respond_to do |format|
        format.xml  { render :layout => false }
      end
    end
  end
end