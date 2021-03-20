document.getElementById("check-availability-button").addEventListener("click", function(){
        
    attention.custom(
      {title:"Choose your dates", 
      html:htmlData,
      willOpen: () => {
            const elem = document.getElementById('reservation-dates-popup');
            const rangepicker = new DateRangePicker(elem, {
               // ...options
               minDate: new Date()
              })
            },
      didOpen: () => {
      document.getElementById('popup-start-date').removeAttribute('disabled');
      document.getElementById('popup-end-date').removeAttribute('disabled');
      },
      callback: function(formValues) {
          let form = document.getElementById("make-reservation-forms");
          console.log(form)
          let formdata = new FormData(form);
          console.log(formValues);
          
          formdata.append("csrf_token", "{{.CSRFToken}}");
          formdata.append("room_id", roomid);
          console.log(formdata);

          console.log("called")
          fetch("/search-json", {method: "post", body: formdata,}).then((response)=> response.json()).then((data)=>{
              if (data.ok) {
                  console.log("room is available!")
                  console.log(data.ok)
                  console.log(data.start_date)
                  attention.custom(
                      {
                          title: "",
                          html: `<a href="/book-room/?id=`+roomid+`&s=`+data.start_date+`&e=`+data.end_date+`">Book Now</a>`,
                          icon: "success",
                          showConfirmButton: false
                      }
                  )
                      
                  
              } else {
                  console.log("sorry not available")
              }
              
          })
      }
  });

    });