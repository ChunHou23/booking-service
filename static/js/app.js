function Prompt() {
    let toast = function(c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            }
        });
        Toast.fire({});
    }

    let success = function(c) {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c;

        Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer
        });
    }

    let error = function(c) {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c;

        Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer
        });
    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showCancelButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showCancelButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
        });
        
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }

        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

function BookRoom(roomId, token) {
    document.getElementById("check-availability-button").addEventListener("click", function() {
        let html = `
            <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
                <div class="form-row">
                    <div class="col">
                        <div class="form-row" id="reservation-dates-modal">
                            <div class="col">
                                <input required type="text" class="form-control" name="start_date" id="start_date" placeholder="Arrival Date">
                            </div>
                            <div class="col">
                                <input required type="text" class="form-control" name="end_date" id="end_date" placeholder="Departure Date">
                            </div>
                        </div>
                    </div>
                </div>
            </form>
        `
        attention.custom({
            msg: html,
            title: "Choose your date",
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal")
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: false,
                    orientation: 'top',
                    minDate: new Date(),
                })
            },
            callback: function(result) {
                console.log("called")

                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", token)
                formData.append("room_id", roomId)

                fetch('search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: "success",
                                msg: '<p>Room is available!</p>'
                                +'<p><a href="/book-room?id='
                                + data.room_id
                                + '&s='
                                + data.start_date
                                + '&e='
                                + data.end_date
                                + '" class="btn btn-primary">'
                                +'Book Now</a></p>',
                                showCancelButton: false,
                            })
                        } else {
                            attention.error({
                                msg: "No Availability"
                            })
                            console.log("room is not available")
                        }
                    })
            }
        })
    })
}
