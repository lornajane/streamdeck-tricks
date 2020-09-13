document.addEventListener('keydown', doKey);

function doKey(e) {
    switch(e.code) {
        case "KeyC":
            a = document.getElementsByClassName('chat-input__textarea');
            a[0].firstChild.focus()
            break;
    }
}
