import { sendWebSocketMessage } from './websocket.js';

export function submitAnswer(playerID, questionID, answer) {
    console.log(`Submitted answer: ${answer}, PlayerID: ${playerID}, QuestionID: ${questionID}`);

    fetch('http://localhost:8080/answer', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            player_id: playerID,
            question_id: questionID,
            answer: answer,
        }),
    })
}




export function showQuestion(question, onSubmitAnswer) {
    const questionElement = document.getElementById('question');
    const choicesElement = document.getElementById('choices');

    questionElement.innerText = question.text;
    choicesElement.innerHTML = '';

    if (question.choices) {
        question.choices.forEach(choice => {
            const button = document.createElement('button');
            button.className = 'choice';
            button.innerText = choice;
            button.onclick = () => onSubmitAnswer(button, choice);
            choicesElement.appendChild(button);
        });
    } else {
        console.error('Question choices are undefined');
    }
}

export function showCorrectAnswer(correctAnswer) {
    const choicesElement = document.getElementById('choices');
    const buttons = document.querySelectorAll('button.choice');
    buttons.forEach(button => {
        if (parseInt(button.innerText) === correctAnswer) {
            button.classList.add('correct');
        } else if (button.classList.contains('selected')) {
            button.classList.add('incorrect');
        }
        button.disabled = true;
    });
    setTimeout(() => {}, 2000);
}
