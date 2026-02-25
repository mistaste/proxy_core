package com.mahsanet.proxy_core

import android.os.Build
import android.util.Log

object ProxyCoreApplication {
    private const val TAG = "ProxyCoreApplication"
    
    init {
        
        try {
            System.loadLibrary("fdsan_workaround")
        } catch (e: UnsatisfiedLinkError) {
            Log.w(TAG, "fdsan_workaround library not found: ${e.message}")
        }
    }
    
    
    private external fun disableFdsanNative()
    
    






    fun initializeFdsanWorkaround() {
        
        
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.Q) {
            try {
                disableFdsanNative()
                Log.i(TAG, "fdsan configured to WARN mode successfully")
            } catch (e: Exception) {
                Log.w(TAG, "Failed to configure fdsan: ${e.message}")
            }
        }
    }
}
