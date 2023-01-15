
# Mail indexer

Mail indexer takes a directory of files and transforms them into semi-structured data and uploads them to zinc search through its own API. 


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`USER_ZINC`

`PASS_ZINC`

`HOST_ZINC`

## Profiling ScreenShots

![Heap memory](https://github.com/lam1391/Mail_indexer/blob/200f2f7da42091f511aca2b53d20def1ad9e7df7/cmd/main/Profiling/profile_allocs_space.png)

![top 5 memory use](https://github.com/lam1391/Mail_indexer/blob/200f2f7da42091f511aca2b53d20def1ad9e7df7/cmd/main/Profiling/profile_top5_funtions.png)

![most consume memory funtion](https://github.com/lam1391/Mail_indexer/blob/200f2f7da42091f511aca2b53d20def1ad9e7df7/cmd/main/Profiling/profile_funtion_most_heavy.png)
