#ifndef PIPELINE_H
#define PIPELINE_H

#include <mutex>

class Pipeline {
public:
    Pipeline (int low, int high);

    void updateThresholds (int low, int high);

    int low ();
    int high ();

private:
    int low_;
    int high_;

    std::mutex mutex;
};

#endif
