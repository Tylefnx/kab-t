// ui.js

export function updateStatus(status) {
    const statusElement = document.getElementById('status');
    statusElement.innerText = status;
}

export function toggleQuizDisplay(show) {
    const quizElement = document.getElementById('quiz');
    quizElement.style.display = show ? 'block' : 'none';
}

export function updateTimer(timeLeft) {
    const timerElement = document.getElementById('timer');
    timerElement.innerText = `Time left: ${timeLeft} seconds`;
}
