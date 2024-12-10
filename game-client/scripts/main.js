document.addEventListener('DOMContentLoaded', () => {
    const statusElement = document.getElementById('status');
    const quizElement = document.getElementById('quiz');
    const questionElement = document.getElementById('question');
    const choicesElement = document.getElementById('choices');
    const timerElement = document.getElementById('timer');
    let timer;
    let currentQuestionID = null;
    let playerID = 'player1'; // Bu ID'yi dinamik olarak belirleyebilirsiniz

    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onopen = () => {
        console.log('Connected to the server');
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log('Received data:', data); // Debugging için
        if (data.status === 'quiz') {
            statusElement.style.display = 'none';
            quizElement.style.display = 'block';
            showQuestion(data.question);
            startTimer(data.timeout);
        } else if (data.status === 'waiting') {
            statusElement.innerText = `Waiting for more players: ${data.count} / ${data.max}`;
        } else if (data.status === 'answer') {
            showCorrectAnswer(data.correct_answer);
        }
    };

    socket.onclose = () => {
        console.log('Disconnected from the server');
    };

    function showQuestion(question) {
        currentQuestionID = question.id;
        questionElement.innerText = question.text;
        choicesElement.innerHTML = '';

        if (question.choices) {
            console.log('Question choices:', question.choices); // Hata ayıklama için
            question.choices.forEach(choice => {
                const button = document.createElement('button');
                button.className = 'choice';
                button.innerText = choice;
                button.onclick = () => submitAnswer(button, choice);
                choicesElement.appendChild(button);
                console.log('Button added:', button); // Hata ayıklama için
            });
        } else {
            console.error('Question choices are undefined');
        }
    }

    function showCorrectAnswer(correctAnswer) {
        const buttons = document.querySelectorAll('button.choice');
        buttons.forEach(button => {
            if (parseInt(button.innerText) === correctAnswer) {
                button.classList.add('correct');
            } else if (button.classList.contains('selected')) {
                button.classList.add('incorrect');
            }
            button.disabled = true; // Sadece bu sorunun butonlarını devre dışı bırakıyoruz
        });
        setTimeout(() => {
            //choicesElement.innerHTML = ''; // Sonraki soru için temizliyoruz
        }, 2000);
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

    function submitAnswer(button, answer) {
        console.log(`Submitted answer: ${answer}`);
        button.classList.add('selected');
        const buttons = document.querySelectorAll('button.choice');
        buttons.forEach(btn => btn.disabled = true); // Şıkları devre dışı bırakıyoruz
        fetch('http://localhost:8080/answer', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                player_id: playerID,
                question_id: currentQuestionID,
                answer: answer,
            }),
        }).then(response => {
            if (response.ok) {
                console.log('Answer submitted successfully');
            } else {
                console.error('Failed to submit answer');
            }
        }).catch(error => {
            console.error('Error submitting answer:', error);
        });
    }
});
