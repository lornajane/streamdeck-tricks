document.addEventListener('keydown', doKey);

function doKey(e) {
    switch(e.code) {
        case "KeyC":
            a = document.getElementsByClassName('messages-container');
            a[0].lastChild.lastChild.focus()
            break;
        case "KeyD":
            i = document.getElementsByClassName('whole-message-container')
            i[0].focus()
            break;
    }
}
