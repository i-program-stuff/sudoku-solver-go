# Sudoku Solver

I have been meaning to learn Golang for a while and thought that making a simple project like a Sudoku solver would be the best way to accomplish that.

### Programming the solver

#### A Problem

While programming the solver, I had to find a way to track which numbers occupied a specific row, column and box. The obvious approach is to use sets, so for this board:

![Sudoku-Board-1](https://github.com/user-attachments/assets/d89948a4-942a-4e8a-83ff-3659673d65ee)

for the upmost rightmost cell, `{5, 9, 4}` occupy the row space while `{9, 4, 7, 6}` occupy the column space and `{4, 9, 7, 3}` occupy the box space. So to find which numbers can go on the upmost rightmost cell you need to find the union of all these sets then find the complement of the union.

```
non_possible_numbers = column_numbers ∪ row_numbers ∪ box_numbers
possible_numbers = 'non_possible_numbers

In our case:
{5, 9, 4} ∪ {9, 4, 7, 6} ∪ {4, 9, 7, 3} = {5, 9, 4, 7, 6, 3}
complement of {5, 9, 4, 7, 6, 3} = {1, 2, 8}
```

so {1, 2, 8} are the possible numbers that can be placed in that cell.

#### OK, so just use sets right?

Well Golang doesn't have sets, but there are ways you can implement it, but being a learning project, I did not want to copy code from the internet, so I started thinking and an idea came to me.

##### Bitmaps as sets

Now I won't claim to be the first person to use this technique, but this was quite a clever idea that I had. If you represent the numbers in the sudoku board like this:

```
let {1} be 1₂
let {2} be 10₂
let {3} be 100₂
so and so...
let {8} be 10000000₂
let {9} be 100000000₂
```

then you can represent a set like `{1, 2, 8}` as 10000011₂.

Operations like union just become the `or` operator and intersection becomes the `and` operator.
The complement becomes `not` (with a mask of 111111111₂ because you have to fit the entire set in a 16bit int).

Now our earlier set example becomes this:
```
{5, 9, 4} ∪ {9, 4, 7, 6} ∪ {4, 9, 7, 3} = {5, 9, 4, 7, 6, 3}
100011000₂ | 101101000₂ | 101001100₂ = 101111100₂

complement of {5, 9, 4, 7, 6, 3} = {1, 2, 8}
~101111100₂ = 10000011₂
```

and voila, the poor man's set. This is also way more efficient as a state can be stored in 2 bytes and all operations are O(1). 

After that you can just use Backtracking to finish the solver.

### Thoughts

This was a pretty fun day long project and as for if I liked Golang? Meh.
