function Prompt() {
    
    let toast = function(c) {
        const {
          title = "",
          position = 'top-end',
          icon = 'success'

        } = c;
        const Toast = Swal.mixin({
        icon : icon,
        title: title,
        toast: true,
        position: position,
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
        didOpen: (toast) => {
          toast.addEventListener('mouseenter', Swal.stopTimer)
          toast.addEventListener('mouseleave', Swal.resumeTimer)
        }
      })

      Toast.fire({})


    }
    let success = function(c) {
      const {
        title="",
        text="",
        footer="",
        icon="success"
      } = c ;

      Swal.fire({
      icon: icon,
      title: title,
      text: text,
      footer: footer,
    })

    }

    let error = function(c) {
      const {
        title="",
        text="",
        footer="",
        icon="error"
      } = c ;

      Swal.fire({
      icon: icon,
      title: title,
      text: text,
      footer: footer,
    })

    }

    let custom = async function(c) {
      const {
        title="",
        html="",
        icon="",
        showConfirmButton=true

      } = c;

      const { value: formValues } = await Swal.fire({
        icon: icon,
        title: title,
        html: html,
        focusConfirm: false,
        showCancelButton: true,
        showConfirmButton: showConfirmButton,
        willOpen: () => {
          if(c.willOpen !== undefined) {
            c.willOpen();
          }
          },

        didOpen: () => {
          if(c.didOpen !== undefined) {
            c.didOpen();
          }
        },
        preConfirm: () => {
          return [
            document.getElementById('popup-start-date').value,
            document.getElementById('popup-end-date').value
          ]
        }
      })

      if (formValues) {
        if (formValues.value !== "") {
          if (c.callback !== undefined)  {
            return c.callback(formValues);
          } 
        } else {
          return c.callback(false);
        }
      }

    }



    return {
      toast : toast,
      success : success,
      error : error,
      custom: custom,
    }
  }

//   var room_id = "1" 

//   function GetBookingResponse(input) {
//     room_id = input
//   }

  