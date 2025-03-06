import React from "react";

import '../css/RulesScreen.css';

function RulesScreen({closeCallback, cards}) {
  return (
    <div className="RulesScreen">
      <div className="RulesScreen-Back-Outer">
        <div className="RulesScreen-Back" onClick={() => closeCallback()}></div>
      </div>
      <div className="RulesScreen-Content">
        <h1>Architorture - Rules</h1>
        <p className="RulesScreen-Text">
          The purpose of this game is to simulate the obstacles, learnings and realities of architecture school. Architorture will help inform current and 
          potential future students understand what tends to be taught and what is actually useful for the future. What can our minds hold onto and what 
          is dropped throughout the years.
          <br/><br/>
          The original edition represents the beginning of architecture school which is 1st year, Bachelors edition. Additional expansions, representing 
          2nd, 3rd and 4th year, could be added if students decide to proceed with an architecture education after playing the original edition.
        </p>
        <h2>Mindset</h2>
        <p className="RulesScreen-Text">
          In this game you are the student and your opponent is architecture school. Players will be playing against other students in competition to see 
          who can best balance the workload and avoid obstacles in architecture school based on their memory capacity and Eureka moments. Players will use 
          knowledge, action and primary cards to aid them along their journey and attempt to reach the end without being overcome by architorture.
        </p>
        <h2>Game Ending</h2>
        <p className="RulesScreen-Text">
          The last player who hasn’t been wiped out by architorture wins the game!
        </p>
        <h2>Gameplay</h2>
        <ol className="RulesScreen-Text">
          <li><p>Play will go in a clockwise direction.</p></li>
          <li><p>Every player is allowed to have 7 cards in hand and 3 cards in the usb archive.</p></li>
          <li><p>Each turn players are allowed to pass or play cards.</p></li>
          <li><p>Players are allowed to play an unlimited amount of action cards per turn.</p></li>
          <li><p>To pass, players must draw a card to end their turn.</p></li>
          <li><p>Players must make sure that they have a maximum of 7 cards in their hand by the end of their turn (including the one that is to be drawn to end their turn). Any extra cards must be played or discarded prior to their turn ending.</p></li>
          <li><p>Players may discard as many cards as they want in hand prior to their turn ending. This may be done so to create more room in their hand, or simply because they feel like there’s no use for the card.</p></li>
          <li><p>Every player ends their turn by drawing a card from the draw pile or playing an action card that ends their turn.</p></li>
        </ol>
        <h2>Discarded Cards</h2>
        <p className="RulesScreen-Text">
          Discarded cards are placed at the bottom of the draw pile.
        </p>
        <h2>Cards in Hand/Memory</h2>
        <ul className="RulesScreen-Text">
          <li><p>The cards in hand represent each players’ brain memory.</p></li>
          <li><p>Players are allowed a default maximum capacity of 7 cards in hand unless expanded by a Memory Expansion Card.</p></li>
        </ul>
        <h2>Cards in USB Archive</h2>
        <ul className="RulesScreen-Text">
          <li><p>Similar to architecture school, most architecture students will have an abundance of USBs and/or hard drives to store their projects and other references.</p></li>
          <li><p>
            A player may archive some cards into their USB for future reference. A player may play 1 action per turn with the USB archive. This includes swapping out a card in memory 
            (in hand) with a card in the usb archive, adding a card to the usb archive, and removing  a card from the usb archive (players must have enough memory/in-hand capacity to do this action).
          </p></li>
          <li><p>Players are allowed a default maximum USB capacity of 3 cards unless expanded by a USB Expansion card.</p></li>
        </ul>
        <h2>Disclaimer</h2>
        <p className="RulesScreen-Text">
          This is just a game, we’re all friends here. Don’t try to sabotage your classmates in real life. The game is based on the experiences at the Azrieli School of Architecture at Carleton University.
        </p>
        <h2>Cards</h2>
        <div className="RulesScreen-Cards">
        {cards.map(function(c, i) {
          return (
            <img className="RulesScreen-Card" src={require("../res/cards/" + c.dbId + "_" + c.expansionId + "_" + c.cardNumber + ".png").default} alt={c.description}/>
          );
        })}
        </div>
      </div>
    </div>
  );
}

export default RulesScreen;
