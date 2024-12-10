document.addEventListener('DOMContentLoaded', () => {
    const statusElement = document.getElementById('status');
    const quizElement = document.getElementById('quiz');
    const questionElement = document.getElementById('question');
    const timerElement = document.getElementById('timer');
    let timer;

    // WebSocket ile sunucuya bağlanma
    const socket = new WebSocket('ws://localhost:8080');

    socket.onopen = () => {
        console.log('Connected to the server');
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.status === 'quiz') {
            statusElement.style.display = 'none';
            quizElement.style.display = 'block';
            showQuestion(data.question);
            startTimer(data.timeout);
        } else if (data.status === 'waiting') {
            statusElement.innerText = `Waiting for more players: ${data.count} / ${data.max}`;
        }
    };

    socket.onclose = () => {
        console.log('Disconnected from the server');
    };

    function showQuestion(question) {
        questionElement.innerText = question.text;
        document.getElementById('answer').value = '';
    }

    function startTimer(seconds) {
        let timeLeft = seconds;
        timerElement.innerText = `Time left: ${timeLeft} seconds`;
        clearInterval(timer);
        timer = setInterval(() => {
            timeLeft -= 1;
            timerElement.innerText = `Time left: ${timeLeft} seconds`;
            if (timeLeft <= 0) {
                clearInterval(timer);
                // Zaman dolduğunda yapılacak işlemler
            }
        }, 1000);
    }
});

function submitAnswer() {
    const answer = document.getElementById('answer').value;
    console.log(`Submitted answer: ${answer}`);
    // Cevabı sunucuya gönderme kodu
}
