xml.instruct! :xml, :version=>"1.0"
xml.encoon do
  @grid.load_cached_grid_structure 
  @grid.row_all(params[:filters], '', params[:page]).each do |row| 
    @grid.row_export(xml, row)
  end
end