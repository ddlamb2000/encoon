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
Encoon::Application.routes.draw do
  root :to => 'grid#home'
  devise_for :users
  match '/*grid_id/__list' => 'grid#list', :as => 'list'
  match '/*grid_id/__new' => 'grid#new', :as => 'new'
  match '/*grid_id/__create' => 'grid#create', :as => 'create', :via => [:post]
  match '/*grid_id/__import' => 'grid#import', :as => 'import'
  match '/*grid_id/__upload' => 'grid#upload', :as => 'upload'
  match '/*grid_id/*id/_details' => 'grid#details', :as => 'details'
  match '/*grid_id/*id/_edit' => 'grid#edit', :as => 'edit'
  match '/*grid_id/*id/_update' => 'grid#update', :as => 'update', :via => [:post]
  match '/*grid_id/*id/_destroy' => 'grid#destroy', :as => 'destroy', :via => [:delete]
  match '/*grid_id/*id/_attach_document' => 'grid#attach_document', :as => 'attach_document'
  match '/*grid_id/*id/_save_attachment' => 'grid#save_attachment', :as => 'save_attachment'
  match '/*grid_id/*id/_delete_attachment' => 'grid#delete_attachment', :as => 'delete_attachment'
  match '/*workspace/*grid_id/*id.xml' => 'grid#export_row', :format => :xml, :as => 'export_row_xml'
  match '/*workspace/*grid_id.xml' => 'grid#export_list', :format => :xml, :as => 'export_list_xml'
  match '/*workspace/*grid_id/*id' => 'grid#show', :as => 'show'
  match '/*workspace/*grid_id' => 'grid#show', :as => 'show'
  match '/__history' => 'grid#history', :as => 'history'
  match '/__set' => 'grid#set', :as => 'set'
  match '/__unset' => 'grid#unset', :as => 'unset'
  match '/__refresh' => 'grid#refresh', :as => 'refresh'
  match '/*workspace' => 'grid#show', :as => 'show'
end