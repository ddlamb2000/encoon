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
class GridController < ApplicationController
  before_filter :load_workspaces, :only => [:home, :show, :refresh]

  before_filter :authenticate_user!, :only => [:create,
                                               :update,
                                               :destroy,
                                               :attach_document,
                                               :save_attachment,
                                               :delete_attachment,
                                               :import,
                                               :upload]

  # Renders the home page using hard-coded references.
  def home
    params[:workspace] = "system"
    params[:grid] = Grid::HOME_GRID_UUID
    params[:row] = Grid::HOME_ROW_UUID
    show
  end

  def history
    session[:show_history] = true
    render :partial => "history"
  end

  # Sets the session flag referenced by the given parameter.
  def set
    session[params[:flag]] = true unless params[:flag].nil?
    render :nothing => true
  end

  # Unsets the session flag referenced by the given parameter.
  def unset
    session[params[:flag]] = false unless params[:flag].nil?
    render :nothing => true
  end

  def refresh
    log_debug "GridController#refresh date=#{params[:home][:session_date]}"
    session[:as_of_date] = Date.strptime(params[:home][:session_date], t('datepicker.decode'))
    redirect_to session[:last_url]
  end

  def show
    log_debug "GridController#show"
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters) if @grid.present?
    selectRow
    @columns = @grid.column_all if @grid.present?
    @grid_cast = Row.select_grid_cast(@grid.uuid, @row.uuid) if @grid.present? and @row.present?
    @attached_grids = Grid.select_referenced_grids(@grid.uuid) if @grid.present?
    set_page_title
    push_history
    render :show, :status => @status
  end

  def export_row
    log_debug "GridController#export_row_xml"
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters) if @grid.present?
    selectRow
  end

  # Renders the content of a list through an Ajax request  
  def list
    log_debug "GridController#list"
    @filters = params[:filters]
    @search = params[:search]
    @page = params[:page]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters) if @grid.present?
    @table_row_count = @grid.row_count(@filters) if @grid.present?
    @table_rows = @grid.row_all(@filters, @search, @page, true) if @grid.present?
    @table_columns = @grid.filtered_columns if @grid.present?
    render :partial => "list"
  end
  
  def export_list
    log_debug "GridController#export_list_xml"
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters) if @grid.present?
    @rows = @grid.row_all(@filters, '', 1, true) if @grid.present?
  end

  # Renders the details of an article through an Ajax request  
  def details
    log_debug "GridController#details"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    if @grid.present? and @row.present?
      @versions = @grid.row_all_versions(@row.uuid)
      @locales = @grid.row_all_locales(@row.uuid, @row.version)
      @audits = @row.all_audits
    end
    session[:show_details] = true
    render :partial => "details"
  end

  # Renders the details of an article through an Ajax request  
  def row
    log_debug "GridController#row"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    @columns = @grid.column_all
    set_page_title
    render :partial => "row"
  end

  def new
    log_debug "GridController#edit"
    @container = params[:container]
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters)
    @filters_uuid = get_filters_uuid(@filters)
    @row = @grid.rows.build
    @row_loc = RowLoc.new
    @grid.row_initialization(@row, @filters)
    render :partial => "edit", :locals => {:new_row => true}
  end

  def edit
    log_debug "GridController#edit"
    @container = params[:container]
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters)
    selectRow
    render :partial => "edit", :locals => {:new_row => false}
  end

  def create
    log_debug "GridController#create"
    saved = false
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters, true)
    @row = @grid.rows.new
    @row.initialization
    begin
      @row.transaction do
        log_debug "GridController#create: initialize transaction"
        @row.begin = param_begin_date
        @row.create_user_uuid = current_user.uuid
        @row.update_user_uuid = current_user.uuid
        log_debug "GridController#create: populate from parameter values"
        populate_from_params
        @grid.load_cached_grid_structure(@filters)
        if @grid.has_translation?
          LANGUAGES.each do |lang, locale|
            log_debug "GridController#create: locale=#{locale}"
            @row_loc = @row.new_loc
            @row_loc.uuid = @row.uuid
            @row_loc.version = @row.version
            @row_loc.base_locale = I18n.locale.to_s
            @row_loc.locale = locale
            @row_loc.name = params[:name]
            @row_loc.description = params[:description]
            log_debug "GridController#create: row_loc_validate"
            if @grid.row_loc_validate(@row, @row_loc, Grid::PHASE_CREATE)
              log_debug "GridController#create: create_row_loc!"
              @grid.create_row_loc!(@row_loc)
            else
              @grid.row_validate(@row, Grid::PHASE_CREATE)
              log_debug "GridController#create: rollback!"
              raise ActiveRecord::Rollback
            end
          end
        end
        log_debug "GridController#create: row_validate"
        if @grid.row_validate(@row, Grid::PHASE_CREATE)
          log_debug "GridController#create: create_row!"
          @grid.create_row!(@row)
          saved = true
        else
          log_debug "GridController#create: rollback!(2)"
          raise ActiveRecord::Rollback
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#create: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#create", invalid
      saved = false
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        format.html { list }
      else
        log_debug "GridController#create: error, @row.errors=#{@row.errors.inspect}"
        format.html { render :json => @row.errors, :status => :unprocessable_entity }
      end
    end
  end
  
  def update
    log_debug "GridController#update"
    saved = false
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGridAndWorkspace
    @grid.load_cached_grid_structure(@filters)
    selectRow
    begin
      @grid.transaction do
        log_debug "GridController#update initialize transaction"
        version = @row.version
        if params[:mode] == Grid::PHASE_NEW_VERSION then
          log_debug "GridController#update new version"
          new_version = @grid.row_max_version(@row.uuid)+1
          @row = @row.clone
          @row.version = new_version
          @row.lock_version = 0
          @row.begin = param_begin_date
          if multiple_versions_ok?
            @row.enabled = params[:enabled]
            @row.create_user_uuid = current_user.uuid
            @row.update_user_uuid = current_user.uuid
            log_debug "GridController#update populate from parameter values"
            populate_from_params
            log_debug "GridController#update row_validate"
            if @grid.row_validate(@row, Grid::PHASE_NEW_VERSION)
              log_debug "GridController#update create_row!"
              @grid.create_row!(@row)
              if @grid.has_translation?
                for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, version)
                  log_debug "GridController#update update row_loc values"
                  @row_loc = @row_loc.clone
                  @row_loc.version = new_version
                  @row_loc.lock_version = 0
                  log_debug "GridController#update locale=#{@row_loc.locale}"
                  if @row_loc.locale == I18n.locale.to_s or 
                     @row_loc.base_locale == I18n.locale.to_s
                    log_debug "GridController#update update row_loc values"
                    @row_loc.base_locale = I18n.locale.to_s
                    @row_loc.name = params[:name]
                    @row_loc.description = params[:description]
                  end
                  log_debug "GridController#update row_loc_validate"
                  if @grid.row_loc_validate(@row, 
                                            @row_loc, 
                                            Grid::PHASE_NEW_VERSION)
                    log_debug "GridController#update create_row_loc!"
                    @grid.create_row_loc!(@row_loc)
                  else
                    log_debug "GridController#update row_validate"
                    @grid.row_validate(@row, Grid::PHASE_UPDATE)
                    log_debug "GridController#update: rollback!"
                    raise ActiveRecord::Rollback
                  end
                end
              end
              saved = true
            else
              raise ActiveRecord::Rollback
            end
          end
        else
          log_debug "GridController#update update existing version"
          @row.begin = param_begin_date
          @row.update_user_uuid = current_user.uuid
          if multiple_versions_ok?
            log_debug "GridController#update versions OK"
            @row.enabled = params[:enabled]
            log_debug "GridController#update populate from parameter values"
            populate_from_params
            log_debug "GridController#update row_validate"
            if @grid.row_validate(@row, Grid::PHASE_UPDATE)
              log_debug "GridController#update update_row!"
              @grid.update_row!(@row)
              if @grid.has_translation?
                for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, version)
                  log_debug "GridController#update locale=#{@row_loc.locale}"
                  if @row_loc.locale == I18n.locale.to_s or 
                     @row_loc.base_locale == I18n.locale.to_s
                    log_debug "GridController#update update row_loc values"
                    @row_loc.base_locale = I18n.locale.to_s
                    @row_loc.name = params[:name]
                    @row_loc.description = params[:description]
                    log_debug "GridController#update row_loc_validate"
                    if @grid.row_loc_validate(@row, 
                                              @row_loc, 
                                              Grid::PHASE_UPDATE)
                      log_debug "GridController#update update_row_loc!"
                      @grid.update_row_loc!(@row_loc)
                    else
                      log_debug "GridController#update: rollback!(3)"
                      raise ActiveRecord::Rollback
                    end
                  end
                end
              end
              saved = true
            else
              log_debug "GridController#update: rollback!(4)"
              raise ActiveRecord::Rollback
            end
          end
        end
      end
      log_debug "GridController#update row_update_dates!"
      @grid.row_update_dates!(@row.uuid)
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#update: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#update", invalid
      saved = false
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        format.html { @refresh_list ? list : row }
      else
        log_debug "GridController#update: error, params=#{params.inspect}"
        format.html { render :json => @row.errors, :status => :unprocessable_entity }
      end
    end
  end

  def destroy
    log_debug "GridController#destroy: params=#{params.inspect}"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    saved = false
    name = @grid.row_title(@row)
    begin
      @row.transaction do
        @row.enabled = false
        if @grid.row_validate(@row, Grid::PHASE_DESTROY)
          @grid.update_row!(@row)
          @grid.destroy_row!(@row)
          @grid.row_update_dates!(@row.uuid)
          if @grid.has_translation?
            for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid)
              @grid.destroy_row_loc!(@row_loc)
            end
          end
          saved = true
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#destroy: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#destroy", invalid
      saved = false
    end
    respond_to do |format|
      if saved
        flash[:notice] = t('transaction.deleted', 
                                :type => @grid, :name => name)
        format.html { redirect_to session[:last_url] }
      else
        log_debug "GridController#destroy: error, params=#{params.inspect}"
        format.html { render :action => :show }
      end
    end
  end

  def attach_document
    log_debug "GridController#attach_document"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    @row_attachment = @row.row_attachments.new
    set_page_title 
  end

  def save_attachment
    log_debug "GridController#save_attachment"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    saved = false
    begin
      @row.transaction do
        @row.remove_attachment!(params[:row_attachment][:attach_document])
        @row_attachment = @row.row_attachments.new
        @row_attachment.update_attributes(params[:row_attachment])
        @row.make_audit(Audit::ATTACH)
        saved = true
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#save_attachment: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#save_attachment", invalid
      saved = false
    end
    respond_to do |format|
      if saved
        name = @grid.row_title(@row)
        flash[:notice] = t('transaction.attached', 
                                :type => @grid, :name => name)
        format.html { redirect_to session[:last_url] }
      else
        log_debug "GridController#save_document: error"
        format.html { render :action => "attach_document" }
      end
    end
  end

  def photo
    log_debug "GridController#photo"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    if params[:photo_id].present?
      @row_attachment = @row.row_attachments.find(params[:photo_id])
    else
      @row_attachment = @row.first_photo
    end
    if @row_attachment.present?
      send_data(@row_attachment.document)
    end
  end

  def file
    log_debug "GridController#file"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    @row_attachment = @row.row_attachments.find(params[:file_id])
    if @row_attachment.present?
      send_data @row_attachment.document, 
                :type => @row_attachment.content_type,
                :filename => @row_attachment.file_name
    end
  end

  def delete_attachment
    log_debug "GridController#delete_attachment"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    selectRow
    saved = false
    begin
      @row.transaction do
        @row_attachment = @row.row_attachments.find(params[:file_id])
        if @row_attachment.present?
          @row_attachment.delete
          @row.make_audit(Audit::DETACH)
        end
      end
      saved = true
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "GridController#delete_attachment: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "GridController#delete_attachment", invalid
      saved = false
    end
    if saved
      flash[:notice] = t('transaction.deleted', 
                              :type => @row_attachment.content_type, 
                              :name => @row_attachment.file_name)
    end
    respond_to do |format|
      format.html { render :action => "show" }
    end
  end

  def import
    log_debug "GridController#import"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
    @page_title = "Import Data"
    @page_icon = "import"
    @upload = Upload.new
  end
  
  def upload
    log_debug "GridController#upload"
    selectGridAndWorkspace
    @grid.load_cached_grid_structure
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

private

  def set_page_title
    @status = 200
    unless @grid.nil?
      unless @row.nil?
        @page_title = t('general.object_name', :type => @grid, :name => @grid.row_title(@row))
        case @grid.uuid
          when Grid::ROOT_UUID then @page_icon = "table"
          when Workspace::ROOT_UUID then @page_icon = "workspace"
          else @page_icon = "entity"
        end
      else
        @status = 404
        @page_title = t('error.no_data')
        @page_icon = "warning"
      end
    else
      @status = 404
      @page_title = t('error.no_grid')
      @page_icon = "warning"
    end
  end

  def populate_from_params
    for column in @grid.column_all
      value = params["row_#{column.physical_column}"]
      log_debug "GridController#populate_from_params" +
                " column.physical_column=#{column.physical_column}" +
                " value=#{value}"
      @row.write_value(column, value)
    end
  end

  # This is used to generate a unique id in the page
  # in order to avoid duplicates.
  # This is based on the concatenation of identifiers used to filter data 
  def get_filters_uuid(filters)
    output = ""
    if filters.present?
      filters.collect do |filter|
         output << filter[:column_uuid]
      end
    end
    output
  end

  def multiple_versions_ok?
    if params[:enabled] == "0" and 
        !@grid.row_enabled_version_exists?(@row, params[:mode] == 'new_version')
      @row.errors.add_to_base(t('error.ena_record'))
      return false
    elsif @grid.row_begin_duplicate_exists?(@row, params[:begin_date])
      @row.errors.add_to_base(t('error.dup_record', :date => params[:begin_date]))
      return false
    end
    true
  end
  
  def selectGridAndWorkspace
    Entity.log_debug "GridController#selectGridAndWorkspace"
    @workspace = nil
    @grid = nil
    if params[:grid].present? and params[:grid] != "0"
      if Entity.uuid?(params[:grid])
        @grid = Grid.select_entity_by_uuid(Grid, params[:grid])
      else
        @grid = Grid.select_entity_by_id(Grid, params[:grid])
      end
      if @grid.nil?
        Entity.log_debug "GridController#selectGridAndWorkspace " + 
                         "Invalid: can't find grid #{params[:grid]}"
      else
        Entity.log_debug "GridController#selectGridAndWorkspace: grid found name=#{@grid.name}"
        @workspace = Workspace.select_entity_by_uuid(Workspace, @grid.workspace_uuid)
        if @workspace.nil?
          Entity.log_debug "GridController#selectGridAndWorkspace " + 
                           "Invalid: can't find workspace #{@grid.workspace_uuid}"
        else
          Entity.log_debug "GridController#selectGridAndWorkspace: workspace found name=#{@workspace.name}"
        end
      end
    end
  end

  def selectRow
    Entity.log_debug "GridController#selectRow"
    @row = nil
    if @grid.present? and params[:row].present? and params[:row] != "0"
      if Entity.uuid?(params[:row])
        @row = @row_loc = @grid.row_select_entity_by_uuid(params[:row])
        unlock_as_of_date
      else
        @row = @row_loc = @grid.row_select_entity_by_id(params[:row])
        lock_as_of_date
      end
      if @row.nil?
        Entity.log_debug "GridController#selectRow " +
                         "Invalid: can't find row with " +
                         "grid_uuid=#{@grid.uuid} " +
                         "and uuid=#{params[:row]}"
      else
        Entity.log_debug "GridController#selectRow: row found name=#{@row.name}"
        if @grid.uuid == Workspace::ROOT_UUID 
          @workspace = Workspace.select_entity_by_uuid(Workspace, @row.uuid)
        elsif @grid.uuid == Grid::ROOT_UUID 
          @workspace = Workspace.select_entity_by_uuid(Workspace, @row.workspace_uuid)
        end
        change_as_of_date(@row)
      end
    end
  end
end