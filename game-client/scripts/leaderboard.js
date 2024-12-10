// leaderboard.js

export function showLeaderboard(scores) {
    const leaderboardElement = document.getElementById('leaderboard');
    const leaderboardList = document.getElementById('leaderboard-list');
    leaderboardElement.style.display = 'block';
    leaderboardList.innerHTML = '';

    scores.sort((a, b) => b.score - a.score);
    scores.forEach((player, index) => {
        const listItem = document.createElement('li');
        listItem.innerText = `${index + 1}. ${player.id} - ${player.score}`;
        leaderboardList.appendChild(listItem);
    });
}

export function hideLeaderboard() {
    const leaderboardElement = document.getElementById('leaderboard');
    leaderboardElement.style.display = 'none';
}
