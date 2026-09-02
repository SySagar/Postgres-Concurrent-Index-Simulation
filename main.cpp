#include <iostream>
#include <vector>
#include <string>
#include <map>

#include <thread>
#include <chrono>
#include <mutex>
#include <condition_variable>

using namespace std;

struct User{
    int  id;
    string name;
};

mutex mtx;
condition_variable cv;

vector<User> table;
map<string,int> idx;   // fake index

vector<User> pending;
bool building = false;
int activeWriters = 0; //  any writers currently running

void synchronize() {

    lock_guard<mutex> lock(mtx);

    cout << "Synchronizing..." << endl;

    for (auto user : pending) {
        idx[user.name] = user.id;
    }

    pending.clear();
}

vector<User> getCurrentSnapShot(vector<User> &table) {
    // Return copy of the table with a lock
    lock_guard<mutex> lock(mtx);

    return table;

}


void buildIndex(){
     building = true;
    vector<User> dbSanpshot;
     dbSanpshot = getCurrentSnapShot(table);

    for(auto user: dbSanpshot){

         cout << "Building index for " << user.name << endl;

        {
        lock_guard<mutex> lock(mtx);
        idx[user.name] = user.id;
        }

          this_thread::sleep_for(chrono::seconds(1)); // simulates millions row transaction scan and indexing - 1 s gap each row
    }

    cout << "Builder waiting..." << endl;

    {
    unique_lock<mutex> lock(mtx);// cv lock reqires unique lock, not lock_guard - so temporarily release the mutex while waiting
 
    // Wait until all writers are done before synchronizing
    cv.wait(lock, [] {
        return activeWriters == 0;
    });

    cout << "Builder woke up!" << endl;
}

    synchronize();

    building = false;
}

void insertUser(int id, string name) {

        {
            // this ensures new insertions are done, with a lock to prevent race
        lock_guard<mutex> lock(mtx);
        activeWriters++;
        }

         cout << "Writer started: " << name << endl;

        this_thread::sleep_for(chrono::seconds(3));

    {
     lock_guard<mutex> lock(mtx);  // lock mutex -> safely execture -> unlocks

     User user{id, name};

    table.push_back(user);

    cout << "Inserted: " << name << endl;

    if (building) {
        pending.push_back(user); // remember it as it wasn't caught during scanning
    }


        activeWriters--;
}
cv.notify_one();
}

int main() {
    
    table.push_back({1, "Alice"});
    table.push_back({2, "Bob"});
    table.push_back({3, "Charlie"});
    table.push_back({3, "Charlie"});
    table.push_back({4, "David"});
    table.push_back({5, "Eve"});

    thread builder(buildIndex);

    this_thread::sleep_for(chrono::seconds(2));

    thread writer(insertUser, 6, "Frank");
     
    writer.join();
    builder.join();


   for (auto entry : idx) {
        cout << entry.first
             << " -> "
             << entry.second
             << endl;
    }

    return 0;
}