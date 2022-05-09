/*
Package game contains the logic for game features such as worlds and entities.

The entire game state is managed by a Game instance. This type has a Run function that is
responsible for updating ('ticking') the world and entities. All game related logic should
execute on the same goroutine that called Run. If a different goroutine needs to execute a
game-related job, it can use the Game.Schedule function to schedule a job to run on the
correct goroutine.
*/
package game
