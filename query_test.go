package main

/* func TestQuery(t *testing.T) {

	db, err := gorm.Open(sqlite.Open("../data/test.db"), &gorm.Config{})
	if err != nil {
		t.Log("failed to connect database")
	}
	date := time.Now().AddDate(0, 0, -4)
	var conferences []data_parser.Conference
	tags := []data_parser.Tag{}
	db.Debug().Preload("Tags", "tag_name IN ?", []string{"javascript", "css", "devops"}).Where("end_date > ?", date).Order("start_date ASC").Find(&conferences)
	var conferences_ []data_parser.Conference
	for _, c:= range conferences {
		fmt.Println(c.Name)
		fmt.Println(c.Tags)
		if len(c.Tags) > 0 {
			conferences_ = append(conferences_, c)
		}
	}
	fmt.Println(len(conferences_))
	for _, t:= range tags {
		//fmt.Println(t.Conferences[0].Name)
		for _, c1:= range t.Conferences {
			for _, t1:= range tags {
				for _, c2:= range t1.Conferences {
					if c1.Name == c2.Name {
						fmt.Print("Duplicate")
					}

				}
				}

			}

	}
}
*/