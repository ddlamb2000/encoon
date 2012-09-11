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
  root :to => 'rows#home'
  devise_for :users
  match '/*grid_id/_list' => 'rows#list', :as => 'list'
  match '/*grid_id/_new' => 'rows#new', :as => 'new'
  match '/*grid_id/_create' => 'rows#create', :as => 'create', :via => [:post]
  match '/*grid_id/_import' => 'rows#import', :as => 'import'
  match '/*grid_id/_upload' => 'rows#upload', :as => 'upload'
  match '/*grid_id/*id/_details' => 'rows#details', :as => 'details'
  match '/*grid_id/*id/_edit' => 'rows#edit', :as => 'edit'
  match '/*grid_id/*id/_update' => 'rows#update', :as => 'update', :via => [:post]
  match '/*grid_id/*id/_destroy' => 'rows#destroy', :as => 'destroy', :via => [:post]
  match '/*grid_id/*id/_attach_document' => 'rows#attach_document', :as => 'attach_document'
  match '/*grid_id/*id/_save_attachment' => 'rows#save_attachment', :as => 'save_attachment'
  match '/*grid_id/*id/_delete_attachment' => 'rows#delete_attachment', :as => 'delete_attachment'
  match '/*workspace/*grid_id.xml' => 'rows#export_list', :format => :xml, :as => 'export_list_xml'
  match '/*workspace/*grid_id/*id.xml' => 'rows#export_row', :format => :xml, :as => 'export_row_xml'
  match '/*workspace/*grid_id/*id' => 'rows#show', :as => 'show'
  match '/history' => 'rows#history', :as => 'history'
  match '/set' => 'rows#set', :as => 'set'
  match '/unset' => 'rows#unset', :as => 'unset'
  match '/refresh' => 'rows#refresh', :as => 'refresh'
end