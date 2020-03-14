package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"testing"
)

type users struct {
	Root  xml.Name `xml:"root"`
	Users []user   `xml:"row"`
}

type user struct {
	ID        int    `xml:"id" json:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age" json:"age"`
	About     string `xml:"about" json:"about"`
	Gender    string `xml:"gender" json:"gender"`
	Name      string `json:"name"`
}

type TestCase struct {
	Result  *SearchResponse
	Request *SearchRequest
	Error   error
}

// Sorting
type byName []user

func (n byName) Len() int           { return len(n) }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

type byAge []user

func (n byAge) Len() int           { return len(n) }
func (n byAge) Less(i, j int) bool { return n[i].Age < n[j].Age }
func (n byAge) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

type byID []user

func (n byID) Len() int           { return len(n) }
func (n byID) Less(i, j int) bool { return n[i].ID < n[j].ID }
func (n byID) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

func TestFindUsers(t *testing.T) {
	tcase := TestCase{
		Request: &SearchRequest{
			Limit:      1,
			Offset:     1,
			Query:      "Boyd Wolf",
			OrderField: "",
			OrderBy:    0,
		},
		Result: &SearchResponse{
			Users: []User{
				{
					Id:     0,
					Name:   "Boyd Wolf",
					Age:    22,
					About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
					Gender: "male",
				},
			},
			NextPage: false,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	sc := &SearchClient{
		URL:         ts.URL,
		AccessToken: "31e5e84005900ee819381d22aa4197ac",
	}

	result, err := sc.FindUsers(*tcase.Request)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if !reflect.DeepEqual(tcase.Result, result) {
		t.Errorf("wrong result, expected %#v, \n got %##v", tcase.Result, result)
	}
}

func TestFindUsersWithNoReqToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		URL: ts.URL,
	}
	tcase := TestCase{
		Request: &SearchRequest{},
	}
	_, err := sc.FindUsers(*tcase.Request)
	if err == nil || err.Error() != "Bad AccessToken" {
		t.Errorf("wrong error or no error")
	}
}

func TestFindUsersWithEmptyQuery(t *testing.T) {
	tcase := &TestCase{
		Request: &SearchRequest{
			Query: "",
		},
		Result: &SearchResponse{
			Users: []User{
				{
					Id:     15,
					Name:   "Allison Valdez",
					Age:    21,
					About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
					Gender: "male",
				},
			},
			NextPage: false,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		URL:         ts.URL,
		AccessToken: "31e5e84005900ee819381d22aa4197ac",
	}
	result, err := sc.FindUsers(*tcase.Request)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if !reflect.DeepEqual(tcase.Result.Users[0], result.Users[0]) {
		t.Errorf("wrong result, expected %#v,\n got %#v", tcase.Result.Users, result.Users[0])
	}

}
func TestFindUsersNotFound(t *testing.T) {
	tcase := &TestCase{
		Request: &SearchRequest{
			Query: "Mark Watney",
		},
		Result: &SearchResponse{
			Users:    []User{},
			NextPage: true,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	sc := &SearchClient{
		URL:         ts.URL,
		AccessToken: "31e5e84005900ee819381d22aa4197ac",
	}
	result, err := sc.FindUsers(*tcase.Request)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	if !reflect.DeepEqual(tcase.Result, result) {
		t.Errorf("wrong result, expected %#v,\n got %#v", tcase.Result, result)
	}

}

func TestFindUsersOrderField(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Request: &SearchRequest{
				Query:      "",
				OrderField: "Age",
			},
			Result: &SearchResponse{
				Users: []User{
					{
						Id:     15,
						Name:   "Allison Valdez",
						Age:    21,
						About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
		},
		TestCase{
			Request: &SearchRequest{
				Query:      "",
				OrderField: "Name",
			},
			Result: &SearchResponse{
				Users: []User{
					{
						Id:     16,
						Name:   "Annie Osborn",
						Age:    35,
						About:  "Consequat fugiat veniam commodo nisi nostrud culpa pariatur. Aliquip velit adipisicing dolor et nostrud. Eu nostrud officia velit eiusmod ullamco duis eiusmod ad non do quis.\n",
						Gender: "female",
					},
				},
				NextPage: false,
			},
		},
		TestCase{
			Request: &SearchRequest{
				Query:      "",
				OrderField: "Id",
			},
			Result: &SearchResponse{
				Users: []User{
					{
						Id:     1,
						Name:   "Hilda Mayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
				},
				NextPage: false,
			},
		},
		TestCase{
			Request: &SearchRequest{
				Query:      "",
				OrderField: "123",
			},
			Result: &SearchResponse{
				Users: []User{
					{
						Id:     16,
						Name:   "Annie Osborn",
						Age:    35,
						About:  "Consequat fugiat veniam commodo nisi nostrud culpa pariatur. Aliquip velit adipisicing dolor et nostrud. Eu nostrud officia velit eiusmod ullamco duis eiusmod ad non do quis.\n",
						Gender: "female",
					},
				},
				NextPage: false,
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for caseNum, item := range cases {

		defer ts.Close()
		sc := &SearchClient{
			URL:         ts.URL,
			AccessToken: "31e5e84005900ee819381d22aa4197ac",
		}
		result, err := sc.FindUsers(*item.Request)
		if err != nil {
			t.Errorf("unexpected error: %#v", err)
		}
		if !reflect.DeepEqual(item.Result, result) {
			t.Errorf("[%d] wrong result, expected %#v,\n got %#v", caseNum, item.Result, result)
		}
	}

}

func setJSONHeadersAndOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("AccessToken")

	if token != "31e5e84005900ee819381d22aa4197ac" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, err := os.Open("./dataset.xml")
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()
	// var userStructForMarshalling users
	var u users
	dataBytes, _ := ioutil.ReadAll(file)
	if err := xml.Unmarshal(dataBytes, &u); err != nil {
		fmt.Printf("error: %v", err)
	}

	for i, v := range u.Users {
		u.Users[i].Name = v.FirstName + " " + v.LastName
	}
	query := r.FormValue("query")
	order := r.FormValue("order_field")
	switch query {

	case "":
		setJSONHeadersAndOK(w)
		sort.Sort(byName(u.Users))
	default:
		setJSONHeadersAndOK(w)
		var resusr []user
		for _, v := range u.Users {
			if query == v.Name {
				resusr = append(resusr, v)
				u.Users = resusr
			}
		}
		if len(resusr) == 0 {
			resusr = append(resusr, u.Users[20])
			u.Users = resusr
		}
	}

	switch order {
	case "Age":
		setJSONHeadersAndOK(w)
		sort.Sort(byAge(u.Users))
	case "Id":
		setJSONHeadersAndOK(w)
		sort.Sort(byID(u.Users))
	case "Name":
		setJSONHeadersAndOK(w)
		sort.Sort(byName(u.Users))
	case "":
		setJSONHeadersAndOK(w)
		sort.Sort(byName(u.Users))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(u.Users)
	if err != nil {
		panic(err)
	}
	w.Write(result)
}
