// Hashcash proof-of-work implementation for browser
class HashcashWorker {
    constructor(bits = 20) {
        this.bits = bits;
    }

    async solve(challenge) {
        return new Promise((resolve) => {
            const parts = challenge.split(':');
            if (parts.length !== 7) {
                resolve(null);
                return;
            }

            let counter = 0;
            const startTime = Date.now();

            const compute = () => {
                const batchSize = 1000;
                
                for (let i = 0; i < batchSize; i++) {
                    const stamp = `${parts[0]}:${parts[1]}:${parts[2]}:${parts[3]}:${parts[4]}:${parts[5]}:${counter.toString(16)}`;
                    
                    if (this.verify(stamp, this.bits)) {
                        const elapsed = (Date.now() - startTime) / 1000;
                        console.log(`Hashcash solved in ${elapsed.toFixed(2)}s after ${counter} attempts`);
                        resolve(stamp);
                        return;
                    }
                    
                    counter++;
                }

                // Continue in next event loop to keep UI responsive
                setTimeout(compute, 0);
            };

            compute();
        });
    }

    verify(stamp, requiredBits) {
        const hash = this.sha256(stamp);
        return this.countLeadingZeroBits(hash) >= requiredBits;
    }

    countLeadingZeroBits(hexStr) {
        let count = 0;
        for (let i = 0; i < hexStr.length; i++) {
            const val = parseInt(hexStr[i], 16);
            if (val === 0) {
                count += 4;
            } else {
                // Count leading zeros in this nibble
                for (let bit = 3; bit >= 0; bit--) {
                    if ((val & (1 << bit)) === 0) {
                        count++;
                    } else {
                        return count;
                    }
                }
                return count;
            }
        }
        return count;
    }

    sha256(str) {
        // Convert string to bytes
        const buffer = new TextEncoder().encode(str);
        
        // Use SubtleCrypto API
        return crypto.subtle.digest('SHA-256', buffer).then(hashBuffer => {
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
            return hashHex;
        });
    }

    // Synchronous version for verification (needed for counting)
    sha256Sync(str) {
        // Simple SHA-256 implementation for verification
        // This is a simplified version - for production, consider using a library
        const buffer = new TextEncoder().encode(str);
        
        // For now, we'll use a workaround with sync crypto
        // In a real implementation, you'd want to use the Web Crypto API properly
        // or a synchronous SHA-256 library
        
        // Fallback: use a simple hash for demo purposes
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            const char = str.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash = hash & hash;
        }
        return Math.abs(hash).toString(16).padStart(64, '0');
    }
}

// Modified to use async crypto properly
HashcashWorker.prototype.verify = function(stamp, requiredBits) {
    return this.sha256(stamp).then(hash => {
        return this.countLeadingZeroBits(hash) >= requiredBits;
    });
};

HashcashWorker.prototype.solve = async function(challenge) {
    const parts = challenge.split(':');
    if (parts.length !== 7) {
        return null;
    }

    let counter = 0;
    const startTime = Date.now();

    while (true) {
        const stamp = `${parts[0]}:${parts[1]}:${parts[2]}:${parts[3]}:${parts[4]}:${parts[5]}:${counter.toString(16)}`;
        const hash = await this.sha256(stamp);
        
        if (this.countLeadingZeroBits(hash) >= this.bits) {
            const elapsed = (Date.now() - startTime) / 1000;
            console.log(`Hashcash solved in ${elapsed.toFixed(2)}s after ${counter} attempts`);
            return stamp;
        }
        
        counter++;
        
        // Yield to UI every 100 attempts
        if (counter % 100 === 0) {
            await new Promise(resolve => setTimeout(resolve, 0));
        }
    }
};

// Export for use in other scripts
window.HashcashWorker = HashcashWorker;
