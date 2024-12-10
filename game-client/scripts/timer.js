let timer;

export function startTimer(seconds, onTickCallback, onTimeoutCallback) {
    let timeLeft = seconds;
    onTickCallback(timeLeft);
    clearInterval(timer);
    timer = setInterval(() => {
        timeLeft -= 1;
        onTickCallback(timeLeft);
        if (timeLeft <= 0) {
            clearInterval(timer);
            onTimeoutCallback();
        }
    }, 1000);
}

export function stopTimer() {
    clearInterval(timer);
}
