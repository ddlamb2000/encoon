xml.instruct! :xml, :version=>"1.0"
xml.encoon do
  if @grid.present?
    if @row.present?
      @grid.row_export(xml, @row)
    end
  end
end