package main

import (
	. "carAvailability/entities"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*

Story: Get non available cars
In order to process a search request
As a carAvailabilty service
I want to get all booked cars and return them

Scenario 1: using car type and date
Given a list of bookings from carBooking service
When I send a GET request to the service with the date and car type I need
Then I should get the list of all cars already booked matching my parameters

Scenario 2: using only car type
Given a list of bookings from carBooking service
When I send a GET request to the service with the car type only
Then I should get the list of all cars with the same type already booked

Scenario 3: using only date
Given a list of bookings from carBooking service
When I send a GET request to the service with the date only
Then I should get the list of all cars already booked at this day

Scenario 4: using no params
Given a list of bookings from carBooking service
When I send a GET request to the service without params
Then I should get [ ]

*/

// https://medium.com/@utkarshmani1997/unit-testing-with-ginkgo-part-2-fe6ed881c635
var _ = Describe("Get non available cars", func() {
	var server *ghttp.Server
	var carBookingServiceMock *ghttp.Server // TODO mock carBooking service
	var type1BookingsSamplesMock string = `[
  {
    "supplier": "Picard",
    "beginBookedDate": "2020-11-05T20:42:16.23145Z",
    "id": 1,
    "arrivalId": 5,
    "arrivalNode": {
      "name": "Avignon-solid",
      "id": 5
    },
    "departureId": 1,
    "departureNode": {
      "name": "Marseille",
      "id": 1
    },
    "carId": 1,
    "car": {
      "id": 1,
      "carTypeId": 1
    },
    "endingBookedDate": "2020-11-05T21:42:16.23145Z"
  },
  {
    "supplier": "Microsoft",
    "beginBookedDate": "2020-11-06T20:42:16.232224Z",
    "id": 3,
    "arrivalId": 5,
    "arrivalNode": {
      "name": "Avignon-solid",
      "id": 5
    },
    "departureId": 4,
    "departureNode": {
      "name": "Paris",
      "id": 4
    },
    "carId": 1,
    "car": {
      "id": 2,
      "carTypeId": 1
    },
    "endingBookedDate": "2020-11-06T21:42:16.232224Z"
  }
]`

	BeforeEach(func() {
		// start a test http server
		server = ghttp.NewServer()
		carBookingServiceMock = ghttp.NewServer()
	})

	AfterEach(func() {
		server.Close()
		carBookingServiceMock.Close()
	})

	Describe("using car type and date", func() {
		When("I send a GET request to the service with the date and car type I need", func() {

			params := url.Values{}
			const expectedCarID = 2

			BeforeEach(func() {
				// Add your handler which has to be called for a given path
				// If there are multiple redirects append all the handlers
				server.AppendHandlers(
					GetNonAvailableCarsRoute,
				)

				carBookingServiceMock.RouteToHandler("GET", "/car-booking/findAll/type/1", func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					io.WriteString(w, type1BookingsSamplesMock)
				})

				params.Add("date", "2020-11-05T20:45:16.23145Z")
				params.Add("carType", "1")
			})

			It("I should get the list of all cars already booked matching my parameters", func() {
				// Simulate incoming GET req
				resp, err := http.Get(server.URL() + "/car-availability/getNonAvailableCars?" + params.Encode())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))

				// Check the response's body
				_, err = ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				Expect(err).ShouldNot(HaveOccurred())

				// Get JSON list
				var carListReturned []Car
				decodeError := json.NewDecoder(resp.Body).Decode(&carListReturned)
				Expect(decodeError).ShouldNot(HaveOccurred())

				//Correct car list returned
				Expect(len(carListReturned)).To(Equal(1))
				Expect(carListReturned[0].Id).To(Equal(expectedCarID))
			})
		})
	})

	//Describe("using car type only", func() {
	//	When("I send a GET request to the service with only the car type I need", func() {
	//
	//		// params := url.Values{}
	//
	//		BeforeEach(func() {
	//			// Add your handler which has to be called for a given path
	//			// If there are multiple redirects append all the handlers
	//			server.AppendHandlers(
	//				GetNonAvailableCarsRoute,
	//			)
	//
	//			//params.Add("carType", "1")
	//		})
	//
	//		It("I should get the list of all cars with the same type already booked", func() {
	//			// Simulate incoming GET req
	//			// resp, err := http.Get(server.URL() + GetBookingRoute)
	//			// TODO correct handling
	//
	//			//Expect(err).ShouldNot(HaveOccurred())
	//			//Expect(resp.StatusCode).Should(Equal(http.StatusOK))
	//			//
	//			//// Check the response's body
	//			//_, err = ioutil.ReadAll(resp.Body)
	//			//defer resp.Body.Close()
	//			//Expect(err).ShouldNot(HaveOccurred())
	//			//
	//			//// Get JSON list
	//			//var carListReturned []Car
	//			//decodeError := json.NewDecoder(resp.Body).Decode(carListReturned)
	//			//Expect(decodeError).ShouldNot(HaveOccurred())
	//
	//			// Correct car list returned
	//			// TODO correct id & number of cars given params
	//			//Expect(len(carListReturned)).To(Equal(1))
	//			//Expect(carListReturned[0].Id).To(Equal(expectedCarID))
	//		})
	//	})
	//})
	//
	//Describe("using date only", func() {
	//	When("I send a GET request to the service with only the date I need", func() {
	//
	//		// params := url.Values{}
	//
	//		BeforeEach(func() {
	//			// Add your handler which has to be called for a given path
	//			// If there are multiple redirects append all the handlers
	//			server.AppendHandlers(
	//				GetNonAvailableCarsRoute,
	//			)
	//
	//			//params.Add("date", "2020-10-18T10:05:25Z")
	//		})
	//
	//		It("should get the list of all cars already booked at this day", func() {
	//			// Simulate incoming GET req
	//			// resp, err := http.Get(server.URL() + GetBookingRoute)
	//			// TODO correct handling
	//
	//			//Expect(err).ShouldNot(HaveOccurred())
	//			//Expect(resp.StatusCode).Should(Equal(http.StatusOK))
	//			//
	//			//// Check the response's body
	//			//_, err = ioutil.ReadAll(resp.Body)
	//			//defer resp.Body.Close()
	//			//Expect(err).ShouldNot(HaveOccurred())
	//			//
	//			//// Get JSON list
	//			//var carListReturned []Car
	//			//decodeError := json.NewDecoder(resp.Body).Decode(carListReturned)
	//			//Expect(decodeError).ShouldNot(HaveOccurred())
	//
	//			// Correct car list returned
	//			// TODO correct id & number of cars given params
	//			//Expect(len(carListReturned)).To(Equal(1))
	//			//Expect(carListReturned[0].Id).To(Equal(expectedCarID))
	//		})
	//	})
	//})
	//
	//Describe("using no params", func() {
	//	When("I send a GET request to the service without any params", func() {
	//
	//		BeforeEach(func() {
	//			// Add your handler which has to be called for a given path
	//			// If there are multiple redirects append all the handlers
	//			server.AppendHandlers(
	//				GetNonAvailableCarsRoute,
	//			)
	//		})
	//
	//		It("I should get ???", func() {
	//			// Simulate incoming GET req
	//			// resp, err := http.Get(server.URL() + GetBookingRoute)
	//			// TODO correct handling
	//
	//			//Expect(err).ShouldNot(HaveOccurred())
	//			//Expect(resp.StatusCode).Should(Equal(http.StatusOK))
	//			//
	//			//// Check the response's body
	//			//_, err = ioutil.ReadAll(resp.Body)
	//			//defer resp.Body.Close()
	//			//Expect(err).ShouldNot(HaveOccurred())
	//			//
	//			//// Get JSON list
	//			//var carListReturned []Car
	//			//decodeError := json.NewDecoder(resp.Body).Decode(carListReturned)
	//			//Expect(decodeError).ShouldNot(HaveOccurred())
	//
	//			// Correct car list returned
	//			// TODO correct id & number of cars given params
	//			//Expect(len(carListReturned)).To(Equal(1))
	//			//Expect(carListReturned[0].Id).To(Equal(expectedCarID))
	//		})
	//	})
	//})

})
