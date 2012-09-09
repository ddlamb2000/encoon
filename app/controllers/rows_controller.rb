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
class RowsController < ApplicationController
  before_filter :load_workspaces, :only => [:home, :show, :refresh]

  before_filter :authenticate_user!, :only => [:create,
                                               :update,
                                               :destrroy,
                                               :attach_document,
                                               :save_attachment,
                                               :delete_attachment,
                                               :import,
                                               :upload]

  # Renders the home page using hard-coded references.
  def home
    params[:grid_id] = Grid::HOME_GRID_UUID
    params[:id] = Grid::HOME_ROW_UUID
    selectGrid
    @grid.load_cached_grid_structure
    selectRow
    set_page_title
    push_history
    render :show, :status => @status
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
    log_debug "RowsController#refresh date=#{params[:home][:session_date]}"
    session[:as_of_date] = Date.strptime(params[:home][:session_date], t('datepicker.decode'))
    redirect_to session[:last_url]
  end

  # Renders the content of a list through an Ajax request  
  def list
    log_debug "RowsController#list"
    @filters = params[:filters]
    @search = params[:search]
    @page = params[:page]
    selectGrid
    @grid.load_cached_grid_structure(@filters)
    @table_row_count = @grid.row_count(@filters)
    @table_rows = @grid.row_all(@filters, @search, @page, true) 
    @table_columns = @grid.filtered_columns
    render :partial => "list"
  end
  
  def show
    log_debug "RowsController#show"
    @filters = params[:filters]
    selectGrid
    @grid.load_cached_grid_structure(@filters)
    selectRow
    if params[:format] == 'xml'
      @rows = @grid.row_all(@filters, '', 1, true) if @row.nil?
    else
      @grid_cast = Row.select_grid_cast(@grid.uuid, @row.uuid) if @grid.present? and @row.present?
      @attached_grids = Grid.select_referenced_grids(@grid.uuid)
      set_page_title
      push_history
      render :show, :status => @status
    end
  end

  # Renders the details of an article through an Ajax request  
  def details
    log_debug "RowsController#details"
    selectGrid
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
    log_debug "RowsController#row"
    selectGrid
    @grid.load_cached_grid_structure
    selectRow
    @columns = @grid.column_all
    set_page_title
    render :partial => "row"
  end

  def edit
    log_debug "RowsController#edit"
    @container = params[:container]
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGrid
    @grid.load_cached_grid_structure(@filters)
    selectRow
    if @row.nil? or @row.uuid.nil?
      @filters_uuid = get_filters_uuid(@filters)
      @row = @grid.rows.build
      @row_loc = RowLoc.new
      @grid.row_initialization(@row, @filters)
      render :partial => "edit", :locals => {:new_row => true}
    else
      render :partial => "edit", :locals => {:new_row => false}
    end
  end

  def create
    log_debug "RowsController#create"
    saved = false
    @filters = params[:filters]
    selectGrid
    @grid.load_cached_grid_structure(@filters, true)
    @row = @grid.rows.new
    @row.initialization
    begin
      @row.transaction do
        @row.begin = param_begin_date
        @row.create_user_uuid = current_user.uuid
        @row.update_user_uuid = current_user.uuid
        populate_from_params
        @grid.load_cached_grid_structure(@filters)
        if @grid.has_translation?
          LANGUAGES.each do |lang, locale|
            @row_loc = @row.new_loc
            @row_loc.uuid = @row.uuid
            @row_loc.version = @row.version
            @row_loc.base_locale = I18n.locale.to_s
            @row_loc.locale = locale
            @row_loc.name = params[:name]
            @row_loc.description = params[:description]
            if @grid.row_loc_validate(@row, @row_loc, Grid::PHASE_CREATE)
              @grid.create_row_loc!(@row_loc)
            else
              @grid.row_validate(@row, Grid::PHASE_CREATE)
              log_debug "RowsController#create: rollback!"
              raise ActiveRecord::Rollback
            end
          end
        end
        if @grid.row_validate(@row, Grid::PHASE_CREATE)
          @grid.create_row!(@row)
          saved = true
        else
          log_debug "RowsController#create: rollback!(2)"
          raise ActiveRecord::Rollback
        end
      end
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowsController#create: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "RowsController#create", invalid
      saved = false
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        format.html { list }
      else
        log_debug "RowsController#create: error, @row.errors=#{@row.errors.inspect}"
        format.html { render :json => @row.errors, :status => :unprocessable_entity }
      end
    end
  end
  
  def update
    log_debug "RowsController#update"
    saved = false
    @refresh_list = params[:refresh_list]
    @filters = params[:filters]
    selectGrid
    @grid.load_cached_grid_structure(@filters)
    selectRow
    begin
      @grid.transaction do
        version = @row.version
        if params[:mode] == Grid::PHASE_NEW_VERSION then
          new_version = @grid.row_max_version(@row.uuid)+1
          @row = @row.clone
          @row.version = new_version
          @row.lock_version = 0
          @row.begin = param_begin_date
          if multiple_versions_ok?
            @row.enabled = params[:row][:enabled]
            @row.create_user_uuid = current_user.uuid
            @row.update_user_uuid = current_user.uuid
            populate_from_params
            if @grid.row_validate(@row, Grid::PHASE_NEW_VERSION)
              @grid.create_row!(@row)
              if @grid.has_translation?
                for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, version)
                  @row_loc = @row_loc.clone
                  @row_loc.version = new_version
                  @row_loc.lock_version = 0
                  if @row_loc.locale == I18n.locale.to_s or 
                     @row_loc.base_locale == I18n.locale.to_s
                    @row_loc.base_locale = I18n.locale.to_s
                    @row_loc.name = params[:name]
                    @row_loc.description = params[:description]
                  end
                  if @grid.row_loc_validate(@row, 
                                            @row_loc, 
                                            Grid::PHASE_NEW_VERSION)
                    @grid.create_row_loc!(@row_loc)
                  else
                    @grid.row_validate(@row, Grid::PHASE_UPDATE)
                    log_debug "RowsController#update: rollback!"
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
          @row.begin = param_begin_date
          @row.update_user_uuid = current_user.uuid
          if multiple_versions_ok?
            @row.enabled = params[:row][:enabled]
            populate_from_params
            if @grid.row_validate(@row, Grid::PHASE_UPDATE)
              @grid.update_row!(@row)
              if @grid.has_translation?
                for @row_loc in @grid.row_loc_select_entity_by_uuid(@row.uuid, 
                                                                    version)
                  if @row_loc.locale == I18n.locale.to_s or 
                     @row_loc.base_locale == I18n.locale.to_s
                    @row_loc.base_locale = I18n.locale.to_s
                    @row_loc.name = params[:name]
                    @row_loc.description = params[:description]
                    if @grid.row_loc_validate(@row, 
                                              @row_loc, 
                                              Grid::PHASE_UPDATE)
                      @grid.update_row_loc!(@row_loc)
                    else
                      log_debug "RowsController#update: rollback!(3)"
                      raise ActiveRecord::Rollback
                    end
                  end
                end
              end
              saved = true
            else
              log_debug "RowsController#update: rollback!(4)"
              raise ActiveRecord::Rollback
            end
          end
        end
      end
      @grid.row_update_dates!(@row.uuid)
    rescue ActiveRecord::RecordInvalid => invalid
      log_debug "RowsController#update: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "RowsController#update", invalid
      saved = false
    end
    change_as_of_date(@row)
    respond_to do |format|
      if saved
        format.html { @refresh_list ? list : row }
      else
        log_debug "RowsController#update: error, params=#{params.inspect}"
        format.html { render :json => @row.errors, :status => :unprocessable_entity }
      end
    end
  end

  def destroy
    log_debug "RowsController#destroy: params=#{params.inspect}"
    selectGrid
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
      log_debug "RowsController#destroy: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "RowsController#destroy", invalid
      saved = false
    end
    respond_to do |format|
      if saved
        flash[:notice] = t('transaction.deleted', 
                                :type => @grid, :name => name)
        format.html { redirect_to session[params[:inline] ? :last_url : :prior_url] }
      else
        log_debug "RowsController#destroy: error, params=#{params.inspect}"
        format.html { render :action => "show" }
      end
    end
  end

  def attach_document
    log_debug "RowsController#attach_document"
    selectGrid
    @grid.load_cached_grid_structure
    selectRow
    @row_attachment = @row.row_attachments.new
    set_page_title 
  end

  def save_attachment
    log_debug "RowsController#save_attachment"
    selectGrid
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
      log_debug "RowsController#save_attachment: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "RowsController#save_attachment", invalid
      saved = false
    end
    respond_to do |format|
      if saved
        name = @grid.row_title(@row)
        flash[:notice] = t('transaction.attached', 
                                :type => @grid, :name => name)
        format.html { redirect_to session[:last_url] }
      else
        log_debug "RowsController#save_document: error"
        format.html { render :action => "attach_document" }
      end
    end
  end

  def photo
    log_debug "RowsController#photo"
    selectGrid
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
    log_debug "RowsController#file"
    selectGrid
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
    log_debug "RowsController#delete_attachment"
    selectGrid
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
      log_debug "RowsController#delete_attachment: invalid=#{invalid.inspect}"
      saved = false
    rescue Exception => invalid
      log_error "RowsController#delete_attachment", invalid
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
    log_debug "RowsController#import"
    selectGrid
    @grid.load_cached_grid_structure
    @page_title = "Import Data"
    @page_icon = "import"
    @upload = Upload.new
  end
  
  def upload
    log_debug "RowsController#upload"
    selectGrid
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
      log_debug "RowsController#populate_from_params" +
                " column.physical_column=#{column.physical_column}" +
                " value=#{params[:row][column.physical_column]}"
      @row.write_value(column,
                       params[:row][column.physical_column])
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
    if params[:row][:enabled] == "0" and 
        !@grid.row_enabled_version_exists?(@row, 
                                           params[:mode] == 'new_version')
      @row.errors.add_to_base(t('error.ena_record'))
      return false
    elsif @grid.row_begin_duplicate_exists?(@row, 
                                            params[:begin_date])
      @row.errors.add_to_base(t('error.dup_record', 
                                     :date => params[:begin_date]))
      return false
    end
    true
  end
  
  def selectGrid
    Entity.log_debug "RowsController#selectGrid"
    @grid = nil
    if Entity.uuid?(params[:grid_id])
      @grid = Grid.select_entity_by_uuid(Grid, params[:grid_id])
    else
      @grid = Grid.select_entity_by_id(Grid, params[:grid_id])
    end
    if @grid.nil?
      Entity.log_debug "RowsController#selectGrid " + 
                       "Invalid: can't find data grid #{params[:grid_id]}"
    else
      Entity.log_debug "RowsController#selectGrid: grid found name=#{@grid.name}"
    end
  end

  def selectRow
    Entity.log_debug "RowsController#selectRow"
    @row = nil
    if params[:id].present? and params[:id] != "0"
      if @grid.present?
        if Entity.uuid?(params[:id])
          @row = @row_loc = @grid.row_select_entity_by_uuid(params[:id])
          unlock_as_of_date
        else
          @row = @row_loc = @grid.row_select_entity_by_id(params[:id])
          lock_as_of_date
        end
        if @row.nil?
          Entity.log_debug "RowsController#selectRow " + 
                           "Invalid: can't find row with " +
                           "grid_id=#{params[:grid_id].to_s}" +
                           " and id=#{params[:id].to_s}"
        else
          Entity.log_debug "RowsController#selectRow: row found name=#{@row.name}"
          change_as_of_date(@row)
        end
      end
    end
  end
end