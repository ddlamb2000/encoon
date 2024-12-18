export const seedData = [
  {
    uuid: 'colors',
    title: 'Colors', 
    cols: [
      {uuid: 'colors-col-1', title: 'Color', type: 'coltypes-row-3'},
      {uuid: 'colors-col-2', title: 'Hex', type: 'coltypes-row-3'},
      {uuid: 'colors-col-3', title: 'Red', type: 'coltypes-row-4'},
      {uuid: 'colors-col-4', title: 'Green', type: 'coltypes-row-4'},
      {uuid: 'colors-col-5', title: 'Blue', type: 'coltypes-row-4'},
    ],
    rows: [
      {uuid: 'colors-row-1', data:['IndianRed', '#CD5C5C', '205', '92', '92']},
      {uuid: 'colors-row-2', data:['LightCoral', '#F08080', '240', '128', '128']},
      {uuid: 'colors-row-2', data:['Salmon', '#FA8072', '250', '128', '114']},
    ]
  }
  ,
  {
    uuid: 'coltypes',
    title: 'Column types',
    cols: [
      {uuid: 'coltypes-col-1', title: 'Type', type: 'coltypes-row-3'},
    ],
    rows: [
      {uuid: 'coltypes-row-1', data:['Any']},
      {uuid: 'coltypes-row-2', data:['Title']},
      {uuid: 'coltypes-row-3', data:['String']},
      {uuid: 'coltypes-row-4', data:['Integer']},
      {uuid: 'coltypes-row-5', data:['Decimal']},
      {uuid: 'coltypes-row-6', data:['Date']},
      {uuid: 'coltypes-row-7', data:['Boolean']},
      {uuid: 'coltypes-row-8', data:['Text']},
      {uuid: 'coltypes-row-9', data:['Grid']},
      {uuid: 'coltypes-row-10', data:['View']},
      {uuid: 'coltypes-row-11', data:['Image']},
      {uuid: 'coltypes-row-12', data:['Video']},
      {uuid: 'coltypes-row-13', data:['Sound']},
    ]
  }
]
