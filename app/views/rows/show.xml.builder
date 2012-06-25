xml.instruct! :xml, :version=>"1.0"
xml.encoon do
  @grid.row_export(xml, @row)
end