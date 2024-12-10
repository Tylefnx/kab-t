import { setupWebSocket } from './websocket.js';
import { startTimer, stopTimer } from './timer.js';
import { submitAnswer, showQuestion, showCorrectAnswer } from './question.js';
import { updateStatus, toggleQuizDisplay, updateTimer } from './ui.js';

document.addEventListener('DOMContentLoaded', () => {
    const playerID = 'player1'; // Bu ID'yi dinamik olarak belirleyebilirsiniz
    let currentQuestionID = null;

    setupWebSocket('ws://localhost:8080/ws', (data) => {
        if (data.status === 'quiz') {
            updateStatus('');
            toggleQuizDisplay(true);
            showQuestion(data.question, (button, answer) => {
                button.classList.add('selected');
                const buttons = document.querySelectorAll('button.choice');
                buttons.forEach(btn => btn.disabled = true);
                submitAnswer(playerID, currentQuestionID, answer);
            });
            currentQuestionID = data.question.id;
            startTimer(data.timeout, updateTimer, () => {
                console.log('Timer finished');
            });
        } else if (data.status === 'waiting') {
            updateStatus(`Waiting for more players: ${data.count} / ${data.max}`);
        } else if (data.status === 'answer') {
            showCorrectAnswer(data.correct_answer);
            stopTimer();
        }
    });
});
