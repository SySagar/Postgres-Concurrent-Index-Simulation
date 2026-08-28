#include <iostream>
#include <vector>
#include <string>
#include <map>

using namespace std;

struct User{
    int  id;
    string name;
};

vector<User> table;
map<string,int> idx;   // fake index

void buildIndex(){
    for(auto user: table){
        idx[user.name] = user.id;
    }
}

int main() {
    
    table.push_back({1, "Alice"});
    table.push_back({2, "Bob"});
    table.push_back({3, "Charlie"});

     buildIndex();


   for (auto entry : idx) {
        cout << entry.first
             << " -> "
             << entry.second
             << endl;
    }

    return 0;
}