function ProcessTabs(id) {
    const current   = document.getElementById(id)
    const tabLinks  = current.querySelectorAll('.tabs a')
    const tabPanels = current.querySelectorAll('.tabs-panel')

    for (let el of tabLinks) {
        el.addEventListener('click', e => {
            e.preventDefault()

            current.querySelector('.tabs li.active').classList.remove('active')
            current.querySelector('.tabs-panel.active').classList.remove('active')

            const parentListItem = el.parentElement
            parentListItem.classList.add('active')

            const index = [...parentListItem.parentElement.children].indexOf(parentListItem)
            const panel = [...tabPanels].filter(el => el.getAttribute('data-index') == index)
                panel[0].classList.add('active')
            }
        )
    }
}

function AnimateNum(id, start, end, duration) {
    let startTimestamp = null

    const step = (timestamp) => {
        if (!startTimestamp) startTimestamp = timestamp

        const progress = Math.min((timestamp - startTimestamp) / duration, 1)

        document.getElementById(id).innerHTML = Math.floor(progress * (end - start) + start)

        if (progress < 1) window.requestAnimationFrame(step)
    }

    window.requestAnimationFrame(step)
}